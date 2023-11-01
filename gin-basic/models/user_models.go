package models

import "gorm.io/gorm"

/*
location: 位置，默认0国内
*/

type UserInfo struct {
	gorm.Model
	Phone    string `json:"phone" gorm:"size:32"`
	Email    string `json:"email" gorm:"size:32"`
	Username string `json:"username" gorm:"size:32"`
	Password string `json:"password" gorm:"size:64"`
	Avatar   string `json:"avatar" gorm:"size:126"`
	Wx       string `json:"wx" gorm:"size:126"`
	Location uint   `json:"location" gorm:"default:0"`
}

func (UserInfo) TableName() string {
	return "user_info"
}
