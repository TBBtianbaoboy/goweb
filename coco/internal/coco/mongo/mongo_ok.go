package coco

import (
	"coco/internal/pkg/mgo"
	"errors"
	"fmt"
	"os"
	myconfig "coco/internal/coco/config"
)

//全局数据库声明
var MongoDB mgo.MgoStruct

var (
	connectMgoError = errors.New("Failed To Connect MongoDB!")
)

func InitMongoDB(){
	err := MongoDB.ConnectMgo(myconfig.ConfigStore.MongoIp,myconfig.ConfigStore.MongoPort,myconfig.ConfigStore.MongoDatabase)
	if err != nil {
		fmt.Println(connectMgoError.Error())
		os.Exit(1)
	}
	fmt.Println("MongoDB Success!")
}

