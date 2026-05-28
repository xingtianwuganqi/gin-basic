package models

type PushDevice struct {
	BaseModel
	UserId      uint   `json:"userId" form:"userId" gorm:"index"`
	DeviceId    string `json:"deviceId" form:"deviceId" gorm:"size:128;uniqueIndex:idx_device_platform"`
	Platform    string `json:"platform" form:"platform" gorm:"size:20;uniqueIndex:idx_device_platform"`
	DeviceToken string `json:"deviceToken" form:"deviceToken" gorm:"size:256"`
	IsActive    bool   `json:"isActive" form:"isActive" gorm:"default:true"`
}

type DeviceTokenRequest struct {
	DeviceToken string `json:"deviceToken" form:"deviceToken" binding:"required"`
	Platform    string `json:"platform" form:"platform"`
}
