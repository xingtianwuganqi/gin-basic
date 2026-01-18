package service

import (
	"errors"
	"gin-basic/models"
	"golang.org/x/crypto/bcrypt"
)

// 模拟用户数据存储
var users = make(map[string]models.User)

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
	}, nil
}