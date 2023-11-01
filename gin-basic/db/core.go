package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"pet-project/models"
	"pet-project/settings"
)

var DB *gorm.DB

func LinkInit() {
	host := settings.Conf.Database.Host
	port := settings.Conf.Database.Port
	database := settings.Conf.Database.DataBase
	username := settings.Conf.Database.Username
	password := settings.Conf.Database.Password
	charset := settings.Conf.Database.Charset
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=Local",
		username,
		password,
		host,
		port,
		database,
		charset)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Error to DB connection ,err" + err.Error())
	}
	autoMigrateTable()
}

func autoMigrateTable() {
	DB.AutoMigrate(&models.UserInfo{})
}
