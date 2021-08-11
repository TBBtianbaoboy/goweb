package coco

import (
	"coco/internal/coco"
	mydb "coco/internal/coco/mongo"
    mylog "coco/internal/coco/log"

	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
)

//插入白名单
func importTable(hip string) iris.Map {
	mes := iris.Map{"code": 100, "message": "failed", "error": ""}
	ip := coco.IpStore{}

	//判断IP是否已经存在
	err := mydb.MongoDB.Find(ip.CollectName(),iris.Map{"_ip_":hip},&ip)
	if err == nil {
		mes["error"] = coco.IpExistError.Error()
		return mes
	}

	//插入
	err = mydb.MongoDB.Insert(ip.CollectName(), iris.Map{"_ip_": hip})
	if err != nil {
		mes["error"] = err.Error()
		return mes
	}
	mylog.LoggerFile.Info("add white table success", zap.String("IP", hip))
	mes = iris.Map{"code": 200, "message": "success"}
	return mes
}

//删除白名单------------logged
func exportTable(hip string) iris.Map {
	mes := iris.Map{"code": 100, "message": "failed", "error": ""}
	ip := coco.IpStore{}
	
	//删除
	err := mydb.MongoDB.Delete(ip.CollectName(), iris.Map{"_ip_": hip})
	if err != nil {
		mes["error"] = coco.IpNotExistError.Error()
		return mes
	}
	mylog.LoggerFile.Info("delete white table success", zap.String("IP", hip))
	mes = iris.Map{"code": 200, "message": "success"}
	return mes
}

//插入白名单处理函数
func ImportTableHandler(ctx iris.Context) {
	ip_h := ctx.FormValue("ip")
	ctx.JSON(importTable(ip_h))
}

//删除白名单处理函数
func ExportTableHandler(ctx iris.Context) {
	ip_h := ctx.FormValue("ip")
	ctx.JSON(exportTable(ip_h))
}
