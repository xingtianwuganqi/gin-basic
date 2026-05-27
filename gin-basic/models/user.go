package models

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty" gorm:"index"`
	Username  string    `json:"username" binding:"required" gorm:"size:50"`
	Email     string    `json:"email" binding:"required,email" gorm:"size:100"`
	Password  string    `json:"password" binding:"required" gorm:"size:255"`
	Phone     string    `json:"phone" gorm:"size:20"`
	Role	  string    `json:"role" gorm:"size:20;default:'user'"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required" validate:"min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required" validate:"min=6,max=100"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

type UserInfoResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}