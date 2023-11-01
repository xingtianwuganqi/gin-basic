package models

import (
	"gorm.io/gorm"
)

type RecordList struct {
	gorm.Model
	User            UserInfo      `json:"user" gorm:"foreignKey:UserId"`
	UserId          uint          `json:"userId"`
	ActionType      PetActionType `json:"action_type" gorm:"foreignKey:PetActionTypeId"`
	PetActionTypeId uint          `json:"petActionTypeId"`
	CustomType      PetCustomType `json:"custom_type" gorm:"foreignKey:PetCustomTypeId"`
	PetCustomTypeId uint          `json:"petCustomTypeId"`
	Spend           float32       `json:"spend"`
	Desc            string        `json:"desc"`
}

type PetActionType struct {
	gorm.Model
	Type       uint   `json:"type" form:"type"`
	ActionName string `json:"action_name" form:"action_name"`
}

func (PetActionType) TableName() string {
	return "pet_action_type"
}

type PetCustomType struct {
	gorm.Model
	User       UserInfo `json:"user" gorm:"foreignKey:UserId"`
	UserId     uint     `json:"userId"`
	CustomName string   `json:"customName" gorm:"size:32"`
}

func (PetCustomType) TableName() string {
	return "pet_custom_type"
}

/*
pet_type : 0:默认值，1：猫咪，2：狗，3：其他
*/

type PetInfo struct {
	gorm.Model
	User     UserInfo `json:"user" gorm:"foreignKey:UserId"`
	UserId   uint     `json:"userId" form:"userId"`
	PetType  uint     `json:"pet_type" form:"pet_type" gorm:"default:0"`
	Age      uint     `json:"age" form:"age" gorm:"default:0"`
	Name     string   `json:"name" form:"name" gorm:"size:32"`
	Avatar   string   `json:"avatar" form:"avatar" gorm:"size:64"`
	BirthDay string   `json:"birthDay" form:"birthDay" gorm:"size:32"`
}

func (PetInfo) TableName() string {
	return "pet_info"
}
