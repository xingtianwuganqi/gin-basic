package models

type AppVersion struct {
	BaseModel
	Platform      string  `json:"platform" form:"platform" gorm:"size:20;uniqueIndex"`
	Version       string  `json:"version" form:"version" gorm:"size:20"`
	MinVersion    *uint   `json:"minVersion" form:"minVersion"`
	LatestVersion uint    `json:"latestVersion" form:"latestVersion"`
	ForceUpdate   int     `json:"forceUpdate" form:"forceUpdate" gorm:"default:0"`
	UpdateTitle   *string `json:"updateTitle" form:"updateTitle" gorm:"size:100"`
	UpdateContent *string `json:"updateContent" form:"updateContent" gorm:"type:text"`
	AppStoreURL   *string `json:"appStoreUrl" form:"appStoreUrl" gorm:"size:255"`
}

type AppVersionRequest struct {
	Platform      string  `json:"platform" form:"platform" binding:"required"`
	Version       string  `json:"version" form:"version" binding:"required"`
	MinVersion    *uint   `json:"minVersion" form:"minVersion"`
	LatestVersion uint    `json:"latestVersion" form:"latestVersion" binding:"required"`
	ForceUpdate   int     `json:"forceUpdate" form:"forceUpdate"`
	UpdateTitle   *string `json:"updateTitle" form:"updateTitle"`
	UpdateContent *string `json:"updateContent" form:"updateContent"`
	AppStoreURL   *string `json:"appStoreUrl" form:"appStoreUrl"`
}

type AppVersionQuery struct {
	Platform string `json:"platform" form:"platform" binding:"required"`
}
