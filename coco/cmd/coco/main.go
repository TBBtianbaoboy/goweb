package main

import (
	"github.com/kataras/iris/v12"

	logRedisMgoInit "coco/internal/coco/init"
	myverify "coco/internal/coco/verify"
	mysign "coco/internal/coco/sign"
	myjwt "coco/internal/coco/jwt"
	myconfig "coco/internal/coco/config"
)

func init(){
	logRedisMgoInit.InitAll()
}

func main() {
	
	app := iris.New()
	
	//获取图片验证码
	app.Get("/verify",myverify.GenerateCaptchaHandler)
	//新建用户注册
	app.Post("/register",mysign.RegisterHandler)
	//登陆验证用户和密码以及验证玛，并获取jwt
	app.Post("/login", mysign.LoginHandler)
	//用户使用token登陆
	app.Get("/{name:string}/home", myjwt.MainFunctionHandler, myjwt.MyAuthenticatedHandler)
	//失效token
	app.Post("/control/dj", myjwt.MakeTokenInvaliHandler)
	//添加用户登录白名单
	app.Post("/control/aw",mysign.ImportTableHandler)
	//删除用户登陆白名单
	app.Post("/control/dw",mysign.ExportTableHandler)
	//注销登陆用户
	app.Get("/logout",mysign.LogoutHandler)

	//定时清理过期jwt
	go myjwt.CleanJwtOnTimeHandler()
	app.Listen(":"+myconfig.GetPort())
}
