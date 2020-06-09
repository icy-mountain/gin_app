package main

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	// データベース
	Dialect = "mysql"
	// ユーザー名
	DBUser = "user1"
	// パスワード
	DBPass = "password"
	// プロトコル
	DBProtocol = "tcp(127.0.0.1:3306)"
	// DB名
	DBName = "go_sample?parseTime=true&loc=Asia%2FTokyo"
)

type User struct {
	gorm.Model
	Name string
}

func db_connect() *gorm.DB {
	connectTemplate := "%s:%s@%s/%s"
	connect := fmt.Sprintf(connectTemplate, DBUser, DBPass, DBProtocol, DBName)
	db, err := gorm.Open(Dialect, connect)
	if err != nil {
		panic("cannnot open DB!")
	} else {
		log.Println("OK!")
	}
	db.AutoMigrate(&User{})
	return db
}

func db_create(db *gorm.DB, user *User) {
	var users []User
	db.Where("name = ?", user.Name).Find(&users)
	if len(users) == 0 {
		log.Println("lets create!")
		db.Create(user)
	}
}
