package coco

import (
	"coco/internal/coco"
	myjwt "coco/internal/coco/jwt"
	mylog "coco/internal/coco/log"
	mydb "coco/internal/coco/mongo"
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo/bson"
	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
)


func LogoutHandler(ctx iris.Context) {
	if err := myjwt.J.CheckJWT(ctx); err != nil {
		ctx.JSON(iris.Map{"code": 100, "message": "failed", "error": coco.ExpiredOrINvalidError.Error()})
		return
	}
	ts := ctx.FormValue("token")

	res := coco.WhiteStore{}
	err := mydb.MongoDB.Find(res.CollectName(),bson.M{"jwt_token":ts},&res)
	if err != nil {
		ctx.JSON(iris.Map{"code": "100", "message": "failed", "error": coco.UserNoLoginError.Error()})
		return
	}
	fmt.Println(res.CollectName())

	token := ctx.Values().Get("jwt").(*jwt.Token)

	foobar := token.Claims.(jwt.MapClaims)
	if foobar["mac"] != myjwt.MacOfMe {
		//---------------------------------------------------------------------------------------------------------非法访问可以写入日志
		err := errors.New(ctx.Host() + "非法访问!!!")
		mylog.LoggerFile.Warn(err.Error())
		ctx.JSON(iris.Map{"code": "100", "message": "failed", "error": err.Error()})
		return
	}
	err = mydb.MongoDB.Delete(res.CollectName(),bson.M{"jwt_token":ts})
	if err != nil {
		mylog.LoggerFile.Error("can not logout,unknow error", zap.String("jwt", ts))
		ctx.JSON(iris.Map{"code": "100", "message": "failed", "error": err.Error()})
		return
	}
	mylog.LoggerFile.Info("logout success", zap.String("jwt", ts))
	ctx.JSON(iris.Map{"code": "200", "message": "success"})
}
