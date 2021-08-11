package coco

import (
	"coco/internal/pkg/config"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type configStr struct {
	Port string `yaml:"port"`
	LogPath string			`yaml:"logpath"`
	LogMaxSize int			`yaml:"logmaxsize"`
	LogMaxBackups int		`yaml:"logmaxbackups"`
	LogMaxAge int			`yaml:"logmaxage"`
	LogCompress bool		`yaml:"logcompress"`
	LogLevel int			`yaml:"loglevel"`

	MongoIp string			`yaml:"mongoip"`
	MongoPort string		`yaml:"mongoport"`
	MongoDatabase string	`yaml:"mongodb"`

	VerifyCodeLong int		`yaml:"vclong"`

	RedisIp string			`yaml:"redisip"`
	RedisPort string		`yaml:"redisport"`
	RedisPasswd string		`yaml:"redispasswd"`
	RedisDB int				`yaml:"redisdb"`

	Secret string			`yaml:"secret"`

	TokenLiveTime int		`yaml:"tokenlivetime"`

	UserPasswdErrLock int	`yaml:"userlocktime"`
}

var ConfigStore = configStr{
	Port: "27358",
	LogPath: "/home/aico/go/coco/logs/coco.log",
	LogMaxSize: 128,
	LogMaxBackups: 30,
	LogMaxAge: 7,
	LogCompress: true,
	LogLevel: 0,
	MongoIp: "127.0.0.1",
	MongoPort: "27017",
	MongoDatabase: "mydb",
	VerifyCodeLong: 4,
	RedisIp: "127.0.0.1",
	RedisPort: "6379",
	RedisPasswd: "",
	RedisDB: 1,
	Secret: "ILoveYou",
	TokenLiveTime: 15,
	UserPasswdErrLock: 5,
}

func ReadConfigFile(filepath string ){
	v,err := config.GetConfigYaml(filepath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	praseYaml(v)
}

func praseYaml(v *viper.Viper){
	err := v.Unmarshal(&ConfigStore)
	if err != nil {
		fmt.Println(errors.New("Parse Config File Error").Error())
		os.Exit(2)
	}
}

func GetPort() string {
	return ConfigStore.Port
}
