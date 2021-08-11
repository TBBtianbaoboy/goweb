package mgo

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type MgoStruct struct {
	session *mgo.Session
	database *mgo.Database
}

//连接mongodb数据库
func (db *MgoStruct)ConnectMgo(ip string,port string,database string) error {
	var err error
	db.session,err = mgo.Dial(ip+":"+port)
	if err != nil {
		return err
	}
	db.database = db.session.DB(database)
	return nil
}


/*用户增删查*/
func (db *MgoStruct)Find(collection string,data bson.M,res interface{}) error {
	return db.database.C(collection).Find(data).One(res)
	// res := User{}
	// err := db.collection.Find(bson.M{"username": user["username"]}).One(&res)
	// err := db.Find(res.CollectName(), bson.M{}, &res)
	// return res, err
}

func (db *MgoStruct)FindAll(collection string,data bson.M,res interface{}) error {
	return db.database.C(collection).Find(data).All(res)
}

func (db *MgoStruct)Insert(collection string,data bson.M) error {
	return db.database.C(collection).Insert(data)
}

func (db *MgoStruct)Delete(collection string,data bson.M) error {
	return db.database.C(collection).Remove(data)
	// err := collection.Remove(iris.Map{"username": user["username"]})
}

/*
jwt白名单增删查
func (db *MgoStruct)FindJwt(data J) (WhiteKit, error) {
	res := WhiteKit{}
	err := db.collection.Find(bson.M{"jwt_token": data.jwttoken}).One(&res)
	return res, err
}

func (db *MgoStruct)InsertJwt(data J) error {
	err := db.collection.Insert(bson.M{"jwt_token": data.jwttoken, "username": data.username})
	return err
}

func (db *MgoStruct)DeleteJwt(data J) error {
	err := db.collection.Remove(bson.M{"jwt_token": data.jwttoken})
	return err
}
*/
