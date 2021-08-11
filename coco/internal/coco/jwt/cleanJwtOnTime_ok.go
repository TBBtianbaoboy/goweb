package coco

import (
	mydb "coco/internal/coco/mongo"
	"coco/internal/coco"
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo/bson"
)

//jwt parser
var jwtpar = new(jwt.Parser)

func cleanJwtOnTime() {
	result := []coco.WhiteStore{}
	err := mydb.MongoDB.FindAll("jwt",nil,&result)
	if err != nil {
		time.Sleep(time.Minute * 5)
		return
	}
	for _, value := range result {
		val := value.Jwt
		parsedToken, _ := jwtpar.Parse(val, J.Config.ValidationKeyGetter)
		if J.Config.Expiration {
			if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
				if expired := claims.VerifyExpiresAt(time.Now().Unix(), true); !expired {
					mydb.MongoDB.Delete("jwt",bson.M{"jwt_token":val})
				}
			}
		}
	}
	time.Sleep(time.Minute * 5)
}

//定时清除过期jwt
func CleanJwtOnTimeHandler() {
	for {
		cleanJwtOnTime()
	}
}
