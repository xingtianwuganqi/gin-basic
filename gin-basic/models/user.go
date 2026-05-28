package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint           `json:"id" form:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt" form:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt" form:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" form:"deletedAt"`
}

type User struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty" gorm:"index"`
	Username  string    `json:"username" binding:"required" gorm:"size:50"`
	Email     string    `json:"email" binding:"required,email" gorm:"size:100"`
	Password  string    `json:"password" binding:"required" gorm:"size:255"`
	Phone     string    `json:"phone" gorm:"size:20"`
	Role      string    `json:"role" gorm:"size:20;default:'user'"`
	DeviceId  string    `json:"deviceId" gorm:"column:device_id;size:128;index"`
	Language  string    `json:"language" form:"language" gorm:"size:20"`
	Region    string    `json:"region" form:"region" gorm:"size:20"`
	Platform  string    `json:"platform" form:"platform" gorm:"size:20"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required" validate:"min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required" validate:"min=6,max=100"`
}

type LoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	DeviceId string `json:"deviceId" form:"deviceId"`
	Language string `json:"language" form:"language"`
	Region   string `json:"region" form:"region"`
	Platform string `json:"platform" form:"platform"`
}

type LoginResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Token    string `json:"token"`
}

type UserInfoResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	DeviceId string `json:"deviceId"`
	Language string `json:"language"`
	Region   string `json:"region"`
	Platform string `json:"platform"`
}

type Page struct {
	PageNum  int `json:"pageNum" form:"pageNum"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

type UserEvent struct {
	BaseModel
	UserId     *uint           `json:"userId" form:"userId" gorm:"index"`
	EventName  string          `json:"eventName" form:"eventName" gorm:"size:64;index;not null"`
	FileId     *string         `json:"fileId,omitempty" form:"fileId" gorm:"size:128;index"`
	Language   string          `json:"language" form:"language" gorm:"size:20;index"`
	Region     string          `json:"region" form:"region" gorm:"size:20;index"`
	Platform   string          `json:"platform" form:"platform" gorm:"size:32;index"`
	AppVersion string          `json:"appVersion" form:"appVersion" gorm:"size:32;index"`
	DeviceId   string          `json:"deviceId" form:"deviceId" gorm:"size:128;index"`
	ClientIP   string          `json:"clientIp" form:"clientIp" gorm:"size:64;index"`
	Extra      json.RawMessage `json:"extra,omitempty" form:"extra" gorm:"type:json"`
}

type UserEventRequest struct {
	EventName  string          `json:"eventName" form:"eventName" binding:"required"`
	FileId     *string         `json:"fileId" form:"fileId"`
	Language   string          `json:"language" form:"language"`
	Region     string          `json:"region" form:"region"`
	Platform   string          `json:"platform" form:"platform"`
	AppVersion string          `json:"appVersion" form:"appVersion"`
	DeviceId   string          `json:"deviceId" form:"deviceId"`
	Extra      json.RawMessage `json:"extra" form:"extra"`
}

type UserEventQuery struct {
	Page
	UserId     *uint  `json:"userId" form:"userId"`
	EventName  string `json:"eventName" form:"eventName"`
	FileId     string `json:"fileId" form:"fileId"`
	Language   string `json:"language" form:"language"`
	Region     string `json:"region" form:"region"`
	Platform   string `json:"platform" form:"platform"`
	AppVersion string `json:"appVersion" form:"appVersion"`
	DeviceId   string `json:"deviceId" form:"deviceId"`
	StartTime  string `json:"startTime" form:"startTime"`
	EndTime    string `json:"endTime" form:"endTime"`
}

type PageResponse struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	PageNum  int         `json:"pageNum"`
	PageSize int         `json:"pageSize"`
}
