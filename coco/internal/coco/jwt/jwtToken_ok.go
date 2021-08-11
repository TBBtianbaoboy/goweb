package coco

import (
	"errors"
	"net"
	"time"
	myconfig "coco/internal/coco/config"
	mydb "coco/internal/coco/mongo"
	mylog "coco/internal/coco/log"
	"coco/internal/coco"

	"github.com/globalsign/mgo/bson"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
)

var (
	jwtInvalidError = errors.New("Jwt Is Invalidation Error!")
)

var mySecret = []byte(myconfig.ConfigStore.Secret)
var MacOfMe string

//利用网卡mac地址来放jwt拷贝
func init() {
	inters, _ := net.Interfaces()
	MacOfMe = inters[1].HardwareAddr.String()
}

//j用在鉴权处理函数检验token
var J = jwt.New(jwt.Config{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	},
	Expiration:    true,
	Extractor:     jwt.FromParameter("token"),
	SigningMethod: jwt.SigningMethodHS256,
})

// 生成token
func GetToken(username string) string {
	now := time.Now() //获取当前时间
	token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"mac": MacOfMe,
		"iat": now.Unix(),
		"exp": now.Add(time.Duration(myconfig.ConfigStore.TokenLiveTime) * time.Minute).Unix(), //token可以存在15分钟
	})

	tokenString, _ := token.SignedString(mySecret)

	err := mydb.MongoDB.Insert("jwt",bson.M{"jwt_token":tokenString,"username":username})
	if err != nil {
		mylog.LoggerFile.Error("Token can't be insert!")
		return "Server Error"
	}
	return tokenString
}

//验证token,防拷贝
func MyAuthenticatedHandler(ctx iris.Context) {
	ts := ctx.FormValue("token")
	if err := J.CheckJWT(ctx); err != nil {
		ctx.JSON(iris.Map{"code": 100, "message": "failed", "error": err.Error()})
		mylog.LoggerFile.Info(err.Error(), zap.String("IP", ctx.Host()), zap.String("jwt", ts))
		return
	}

	res := coco.WhiteStore{}
	err := mydb.MongoDB.Find(res.CollectName(),bson.M{"jwt_token":ts},&res)
	if err != nil {
		ctx.JSON(iris.Map{"code": "100", "message": "failed", "error": jwtInvalidError.Error()})
		return
	}

	name := ctx.Values().Get("name")
	if name != res.Username {
		mylog.LoggerFile.Info(res.Username+"!="+name.(string), zap.String("IP", ctx.Host()), zap.String("jwt", ts))
		ctx.JSON(iris.Map{"code": 100, "message": "failed", "error": "Out of Privilege Error!"})
		return
	}

	token := ctx.Values().Get("jwt").(*jwt.Token)

	foobar := token.Claims.(jwt.MapClaims)
	if foobar["mac"] != MacOfMe {
		err := errors.New(ctx.Host() + "非法访问!!!")
		mylog.LoggerFile.Error(err.Error(), zap.String("source mac", foobar["mac"].(string)))
		ctx.JSON(iris.Map{"code": "100", "message": "failed", "error": err.Error()})
		return
	}
	ctx.JSON(iris.Map{"code": "200", "message": "success"})
}

//用户登陆后的第一处理函数
func MainFunctionHandler(ctx iris.Context) {
	username := ctx.Params().Get("name")
	_, b := ctx.Values().Set("name", username)
	if b {
		ctx.Next()
	} else {
		ctx.JSON(iris.Map{"code": 100, "message": "failed", "error": "Unknow Error"})
	}
}

//让token主动失效
func MakeTokenInvaliHandler(ctx iris.Context) {
	mes := iris.Map{"code": "100", "message": "failed", "error": ""}
	jwt_string := ctx.FormValue("jwt")
	res := coco.WhiteStore{}
	
	err := mydb.MongoDB.Find(res.CollectName(),bson.M{"jwt_token":jwt_string},&res)
	if err != nil {
		mes["error"] = jwtInvalidError.Error()
		ctx.JSON(mes)
		return
	}

	err = mydb.MongoDB.Delete(res.CollectName(),bson.M{"username":res.Username})
	if err != nil {
		mylog.LoggerFile.Error("Can't delete exist jwt", zap.String("jwt", jwt_string))
		mes["error"] = "Can't Delete Exist jwt"
		ctx.JSON(mes)
		return
	}
	ctx.JSON(iris.Map{"code": 200, "message": "success"})
}
