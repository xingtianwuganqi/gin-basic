package service

import (
	"errors"
	"gin-basic/db"
	"gin-basic/middleware"
	"gin-basic/models"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 模拟用户数据存储
var users = make(map[string]models.User)

var ErrInvalidQueryTime = errors.New("invalid time format")

// HashPassword 对密码进行哈希处理
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword 检查密码是否正确
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// RegisterUser 用户注册
func RegisterUser(req models.RegisterRequest) error {
	// 检查用户是否已存在
	for _, user := range users {
		if user.Email == req.Email {
			return errors.New("user already exists")
		}
	}

	// 哈希密码
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return err
	}

	// 创建新用户
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	users[user.Email] = user
	return nil
}

// LoginUser 用户登录
func LoginUser(req models.LoginRequest) (models.LoginResponse, error) {
	if strings.TrimSpace(req.DeviceId) != "" {
		return LoginUserByDevice(req)
	}

	if req.Email == "" || req.Password == "" {
		return models.LoginResponse{}, errors.New("email and password are required")
	}

	user, exists := users[req.Email]
	if !exists {
		return models.LoginResponse{}, errors.New("user not found")
	}

	if !CheckPassword(req.Password, user.Password) {
		return models.LoginResponse{}, errors.New("incorrect password")
	}

	// 生成模拟token
	token := "mock-token-for-" + user.Email

	return models.LoginResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Token:    token,
	}, nil
}

func LoginUserByDevice(req models.LoginRequest) (models.LoginResponse, error) {
	deviceID := strings.TrimSpace(req.DeviceId)
	if deviceID == "" {
		return models.LoginResponse{}, errors.New("deviceId is required")
	}

	user := models.User{}
	result := db.DB.Model(&models.User{}).Where("device_id = ?", deviceID).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		user = models.User{
			DeviceId: deviceID,
			Language: strings.TrimSpace(req.Language),
			Region:   strings.TrimSpace(req.Region),
			Platform: strings.TrimSpace(req.Platform),
		}
		if err := db.DB.Create(&user).Error; err != nil {
			return models.LoginResponse{}, err
		}
	} else if result.Error != nil {
		return models.LoginResponse{}, result.Error
	} else {
		updates := map[string]interface{}{}
		if strings.TrimSpace(req.Language) != "" {
			updates["language"] = strings.TrimSpace(req.Language)
			user.Language = strings.TrimSpace(req.Language)
		}
		if strings.TrimSpace(req.Region) != "" {
			updates["region"] = strings.TrimSpace(req.Region)
			user.Region = strings.TrimSpace(req.Region)
		}
		if strings.TrimSpace(req.Platform) != "" {
			updates["platform"] = strings.TrimSpace(req.Platform)
			user.Platform = strings.TrimSpace(req.Platform)
		}
		if len(updates) > 0 {
			if err := db.DB.Model(&user).Updates(updates).Error; err != nil {
				return models.LoginResponse{}, err
			}
		}
	}

	token, err := middleware.GenToken(user.DeviceId)
	if err != nil {
		return models.LoginResponse{}, err
	}

	return models.LoginResponse{
		ID:    user.ID,
		Token: token,
	}, nil
}

// GetUserInfo 获取用户信息
func GetUserInfo(email string) (models.UserInfoResponse, error) {
	user, exists := users[email]
	if !exists {
		return models.UserInfoResponse{}, errors.New("user not found")
	}

	return models.UserInfoResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
		DeviceId: user.DeviceId,
		Language: user.Language,
		Region:   user.Region,
		Platform: user.Platform,
	}, nil
}

func UpdateDeviceToken(userId uint, req models.DeviceTokenRequest) (*models.PushDevice, error) {
	user := models.User{}
	if err := db.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, err
	}

	platform := strings.TrimSpace(req.Platform)
	if platform == "" {
		platform = strings.TrimSpace(user.Platform)
	}
	if platform == "" {
		platform = "ios"
	}

	pushDevice := &models.PushDevice{}
	err := db.DB.Where("device_id = ? and platform = ?", user.DeviceId, platform).First(pushDevice).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	pushDevice.UserId = userId
	pushDevice.DeviceId = user.DeviceId
	pushDevice.Platform = platform
	pushDevice.DeviceToken = strings.TrimSpace(req.DeviceToken)
	pushDevice.IsActive = true

	if pushDevice.ID == 0 {
		if err := db.DB.Create(pushDevice).Error; err != nil {
			return nil, err
		}
		return pushDevice, nil
	}

	if err := db.DB.Save(pushDevice).Error; err != nil {
		return nil, err
	}
	return pushDevice, nil
}

func CreateUserEvent(userId *uint, req models.UserEventRequest, clientIP string) (*models.UserEvent, error) {
	userEvent := &models.UserEvent{
		UserId:     userId,
		EventName:  req.EventName,
		FileId:     req.FileId,
		Language:   req.Language,
		Region:     req.Region,
		Platform:   req.Platform,
		AppVersion: req.AppVersion,
		DeviceId:   req.DeviceId,
		ClientIP:   clientIP,
		Extra:      req.Extra,
	}
	if err := db.DB.Create(userEvent).Error; err != nil {
		return nil, err
	}
	return userEvent, nil
}

func ListUserEvent(query models.UserEventQuery) ([]models.UserEvent, int64, error) {
	if query.PageSize <= 0 {
		query.PageSize = 10
	}
	if query.PageNum <= 0 {
		query.PageNum = 1
	}

	dbQuery := db.DB.Model(&models.UserEvent{})

	if query.UserId != nil {
		dbQuery = dbQuery.Where("user_id = ?", *query.UserId)
	}
	if query.EventName != "" {
		dbQuery = dbQuery.Where("event_name = ?", query.EventName)
	}
	if query.FileId != "" {
		dbQuery = dbQuery.Where("file_id = ?", query.FileId)
	}
	if query.Language != "" {
		dbQuery = dbQuery.Where("language = ?", query.Language)
	}
	if query.Region != "" {
		dbQuery = dbQuery.Where("region = ?", query.Region)
	}
	if query.Platform != "" {
		dbQuery = dbQuery.Where("platform = ?", query.Platform)
	}
	if query.AppVersion != "" {
		dbQuery = dbQuery.Where("app_version = ?", query.AppVersion)
	}
	if query.DeviceId != "" {
		dbQuery = dbQuery.Where("device_id = ?", query.DeviceId)
	}

	startTime, err := parseQueryTime(query.StartTime)
	if err != nil {
		return nil, 0, err
	}
	if startTime != nil {
		dbQuery = dbQuery.Where("created_at >= ?", *startTime)
	}

	endTime, err := parseQueryTime(query.EndTime)
	if err != nil {
		return nil, 0, err
	}
	if endTime != nil {
		dbQuery = dbQuery.Where("created_at <= ?", *endTime)
	}

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	eventList := make([]models.UserEvent, 0)
	if err := dbQuery.
		Order("created_at desc").
		Limit(query.PageSize).
		Offset((query.PageNum - 1) * query.PageSize).
		Find(&eventList).Error; err != nil {
		return nil, 0, err
	}
	return eventList, total, nil
}

func parseQueryTime(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}

	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, layout := range layouts {
		parsed, err := time.ParseInLocation(layout, value, time.Local)
		if err == nil {
			return &parsed, nil
		}
	}

	return nil, ErrInvalidQueryTime
}
