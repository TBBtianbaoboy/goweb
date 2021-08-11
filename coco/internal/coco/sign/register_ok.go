package coco

import (
	"coco/internal/coco"
	mylog "coco/internal/coco/log"
	mydb "coco/internal/coco/mongo"
	myredis "coco/internal/coco/redis"
	"github.com/kataras/iris/v12"
)

func registerUser(nu iris.Map) iris.Map {
	mes := iris.Map{"code": 0, "err": "", "msg": "failed"}
	u := coco.UserStore{}

	//验证用户是否已存在
	err := mydb.MongoDB.Find(u.CollectName(), iris.Map{"username": nu["username"]}, &u)
	if err == nil {
		mes["err"] = coco.UsernameRepeatError.Error()
		return mes
	}

	//验证是否可以新建用户
	if err = mydb.MongoDB.Insert(u.CollectName(), nu); err != nil {
		mes["err"] = coco.CreateAccountError.Error()
		return mes
	}

	//将注册成功信息写入redis
	err = myredis.AddUserToRedis("login", nu["username"], nu["password"])
	if err != nil {
		mylog.LoggerFile.Error("cann't add register user to redis cache")
	}

	mes = iris.Map{"code": 200, "username": nu["username"], "msg": "success"}
	return mes
}

func RegisterHandler(ctx iris.Context) {
	us := iris.Map{"username": ctx.FormValue("username"), "password": ctx.FormValue("password")}
	ctx.JSON(registerUser(us))
}
