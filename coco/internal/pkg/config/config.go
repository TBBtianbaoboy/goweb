package config

import (
	"errors"
	"github.com/spf13/viper"
)

var (
	readConfigFileError = errors.New("Read Config File Error!")
)

func GetConfigYaml(filePath string) (*viper.Viper,error) {
	v := viper.New()
	
	//设置配置文件的文件名,文件类型,和路径
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(filePath)

	if err:= v.ReadInConfig();err != nil {
		return nil,readConfigFileError
	}
	return v,nil
}
