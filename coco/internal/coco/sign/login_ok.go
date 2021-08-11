package coco

import (
	"coco/internal/coco"
	myjwt "coco/internal/coco/jwt"
	mylog "coco/internal/coco/log"
	mydb "coco/internal/coco/mongo"
	myredis "coco/internal/coco/redis"
	myverify "coco/internal/coco/verify"
	"net"

	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
)

//登陆检验
func loginFunc(user iris.Map, ver coco.V, ipString string) iris.Map {
	mes := iris.Map{"code": 0, "err": "", "msg": "failed", "verify": false}
	// ip := coco.IpStore{Collection: "login"}
	u := coco.UserStore{}

	//判断是否在白名单中
	// err := mydb.MongoDB.Find(ip.CollectName(), iris.Map{"_ip_": ipString}, &ip)
	// if err != nil {
		// mylog.LoggerFile.Warn("no privilege user try to connect server", zap.String("IP", ipString))
		// mes["err"] = coco.PrivilegeError.Error()
		// return mes
	// }

	username_1 := user["username"].(string)
	password_1 := user["password"].(string)
	
	flag, err := myredis.FindUserInRedis("login", username_1, password_1)
	if err == coco.ForbiddenLoginError {
		mes["err"] = err.Error()
		return mes
	} else if err == coco.UserHasLoginError {
		mes["err"] = err.Error()
		return mes
	} else if err == coco.PasswordError {
		//检测是否封停账号
		myredis.CheckPasswdRedis("forbidden", username_1)
		mes["err"] = err.Error()
		return mes
	} else if !flag { //redis中不存在，去mongo找
		//判断账号密码是否正确
		err = mydb.MongoDB.Find(u.CollectName(),iris.Map{"username":username_1},&u)
		if err != nil {
			mes["err"] = coco.UsernameError.Error()
			return mes
		}
		//判断用户是否已经登陆
		err = myredis.CheckUserHasLogin(iris.Map{"username":username_1})
		if err != nil {
			mes["err"] = err.Error()
			return mes
		}
		//判断用户是否被禁止登陆
		if myredis.CheckLoginRedis(username_1) {
			mes["err"] = coco.ForbiddenLoginError.Error()
			return mes
		} else if u.Passwd != user["password"] {
			myredis.CheckPasswdRedis("forbidden", username_1)
			mes["err"] = coco.PasswordError.Error()
			return mes
		}
	} //else 表示账号密码都正确

	//判断验证码是否正确
	if ver["capId"] != "" && ver["bas64"] != "" && myverify.CaptchaVerifyHandle(ver) {
		tn := myjwt.GetToken(username_1)
		mylog.LoggerFile.Info("login success", zap.String("username", username_1), zap.String("jwt", tn))
		mes = iris.Map{"code": 1, "token": tn, "msg": "success", "verify": true}
		//登陆成功后清理redis错误记录
		myredis.CleanForbidden("forbidden", username_1)
	} else {
		mes["err"] = coco.VerifyError.Error()
	}
	return mes
}

//登陆处理函数
func LoginHandler(c iris.Context) {
	mes := iris.Map{"username": c.FormValue("username"), "password": c.FormValue("password")}
	ver := coco.V{"capId": c.FormValue("capId"), "bas64": c.FormValue("bas64")}
	ip,_,_ := net.SplitHostPort(c.Host())
	c.JSON(loginFunc(mes, ver, ip))
}
