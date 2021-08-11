package coco

import (
	"coco/internal/coco"
	myconfig "coco/internal/coco/config"

	"github.com/kataras/iris/v12"
	"github.com/mojocn/base64Captcha"
)

var store = base64Captcha.DefaultMemStore

//生成base64验证码-----------------/verify
func GenerateCaptchaHandler(ctx iris.Context){
	drive := base64Captcha.NewDriverDigit(80, 240, myconfig.ConfigStore.VerifyCodeLong, 0.7, 80)
	ca := base64Captcha.NewCaptcha(drive, store)
	id, b64s, err := ca.Generate()
	body := iris.Map{"code": 200, "data": b64s, "captchaId": id, "msg": "success"}
	if err != nil {
		body = iris.Map{"code": 100, "msg": err.Error()}
	}
	ctx.JSON(body)
}

//辅助函数-----loginAndVerify
func CaptchaVerifyHandle(vc coco.V) bool {

	id := vc["capId"]
	bas64 := vc["bas64"]
	if store.Verify(id, bas64, true) {
		return true
	}
	return false
}
