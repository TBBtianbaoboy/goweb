package init

import (
	mylog "coco/internal/coco/log"
	myconfig "coco/internal/coco/config"
	mydb "coco/internal/coco/mongo"
	myredis "coco/internal/coco/redis"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func InitAll(){
	pflag.String("config","/go/coco/configs","config.yaml path")
	viper.BindPFlags(pflag.CommandLine)
	//读取配置文件并反序列化到全局结构体变量
	myconfig.ReadConfigFile(viper.GetString("config"))	
	mylog.InitLogger()
	mydb.InitMongoDB()
	myredis.InitRedis()
}
