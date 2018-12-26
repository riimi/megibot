package model

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
)

var DB *gorm.DB

func InitDB(source string) *gorm.DB {
	db, err := gorm.Open("mysql", source)
	if err != nil {
		log.Fatal(err)
	}
	DB = db
	return db
}
