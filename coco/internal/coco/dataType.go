package coco

import (
	"errors"

	"github.com/globalsign/mgo/bson"
)

//接受数据类型，改进时用结构体替换
type V map[string]string

//用户名单
type UserStore struct {
	Id         bson.ObjectId `bson:"_id"`
	Name       string        `bson:"username"`
	Passwd     string        `bson:"password"`
} //全部写完之后看看是否可以用interface封装这些重复的方法

func (u *UserStore) CollectName() string {
	return "user"
}

//jwt白名单
type WhiteStore struct {
	Id         bson.ObjectId `bson:"_id"`
	Jwt        string        `bson:"jwt_token"`
	Username   string        `bson:"username"`
}

func (u *WhiteStore) CollectName() string {
	return  "jwt"
}

//IP白名单
type IpStore struct {
	Ip_        string `bson:"_ip_"`
}

func (u *IpStore) CollectName() string {
	return "login"
}

//错误类型声明
var (
	DataBaseError       = errors.New("Connect Database Error!")
	UsernameError       = errors.New("Username Error!")
	PasswordError       = errors.New("Password Error!")
	PrivilegeError      = errors.New("No Privilege Error!")
	VerifyError         = errors.New("Verify Code Error!")
	ForbiddenLoginError = errors.New("Forbidden Login!")
	ConnectRedisError = errors.New("Connect Redis Error!")
	AddRedisMapError  = errors.New("Insert Redis Error!")
	UserHasLoginError = errors.New("User Has Login Error!")
	IpExistError    = errors.New("Ip Existed Error!")
	IpNotExistError = errors.New("Ip Not Existed Error!")
	UserNoLoginError      = errors.New("User No Login Error!")
	ExpiredOrINvalidError = errors.New("Token Is Invalid Or Expired!")
	UsernameRepeatError = errors.New("Username Repeated Error!")
	CreateAccountError  = errors.New("Create New Account Error!")
)

