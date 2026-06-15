package biz

import (
	"context"
	"gorm.io/gorm"
	"time"
)

// User 用户实体
type User struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Username  string         `gorm:"unique;not null"`
	Email     string         `gorm:"unique;not null"`
	Password  string         `gorm:"not null"`
}

// UserRepo 用户仓库接口
type UserRepo interface {
	Create(ctx context.Context, user *User) error
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByID(ctx context.Context, id int64) (*User, error)
}

// UserUsecase 用户用例
type UserUsecase struct {
	repo UserRepo
}

// NewUserUsecase 创建用户用例
func NewUserUsecase(repo UserRepo) *UserUsecase {
	return &UserUsecase{repo: repo}
}

// Register 用户注册
func (uc *UserUsecase) Register(ctx context.Context, username, email, password string) (*User, error) {
	user := &User{
		Username: username,
		Email:    email,
		Password: password,
	}
	if err := uc.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// Login 用户登录
func (uc *UserUsecase) Login(ctx context.Context, username, password string) (*User, error) {
	return uc.repo.FindByUsername(ctx, username)
}

// GetUserByID 根据ID获取用户
func (uc *UserUsecase) GetUserByID(ctx context.Context, id int64) (*User, error) {
	return uc.repo.FindByID(ctx, id)
}
