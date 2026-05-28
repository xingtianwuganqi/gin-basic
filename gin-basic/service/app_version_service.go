package service

import (
	"gin-basic/db"
	"gin-basic/models"

	"gorm.io/gorm"
)

func UpsertAppVersion(req models.AppVersionRequest) (*models.AppVersion, error) {
	appVersion := &models.AppVersion{}
	err := db.DB.Where("platform = ?", req.Platform).First(appVersion).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	appVersion.Platform = req.Platform
	appVersion.Version = req.Version
	appVersion.MinVersion = req.MinVersion
	appVersion.LatestVersion = req.LatestVersion
	appVersion.ForceUpdate = req.ForceUpdate
	appVersion.UpdateTitle = req.UpdateTitle
	appVersion.UpdateContent = req.UpdateContent
	appVersion.AppStoreURL = req.AppStoreURL

	if appVersion.ID == 0 {
		if err := db.DB.Create(appVersion).Error; err != nil {
			return nil, err
		}
		return appVersion, nil
	}

	if err := db.DB.Save(appVersion).Error; err != nil {
		return nil, err
	}
	return appVersion, nil
}

func ListAppVersion() ([]models.AppVersion, error) {
	versionList := make([]models.AppVersion, 0)
	if err := db.DB.Order("updated_at desc").Find(&versionList).Error; err != nil {
		return nil, err
	}
	return versionList, nil
}

func GetAppVersion(platform string) (*models.AppVersion, error) {
	appVersion := &models.AppVersion{}
	if err := db.DB.Where("platform = ?", platform).First(appVersion).Error; err != nil {
		return nil, err
	}
	return appVersion, nil
}
