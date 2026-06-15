package data

import (
	"context"

	"my_knowledges/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserRepo 创建用户仓库
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// Create 创建用户（密码自动加密）
func (r *userRepo) Create(ctx context.Context, user *biz.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		r.log.Errorf("Failed to hash password: %v", err)
		return err
	}
	user.Password = string(hashedPassword)

	result := r.data.MySQL.WithContext(ctx).Create(user)
	if result.Error != nil {
		r.log.Errorf("Failed to create user: %v", result.Error)
		return result.Error
	}
	return nil
}

// FindByUsername 根据用户名查找用户
func (r *userRepo) FindByUsername(ctx context.Context, username string) (*biz.User, error) {
	var user biz.User
	result := r.data.MySQL.WithContext(ctx).Where("username = ? AND deleted_at IS NULL", username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.log.Errorf("Failed to find user by username: %v", result.Error)
		return nil, result.Error
	}
	return &user, nil
}

// FindByID 根据ID查找用户
func (r *userRepo) FindByID(ctx context.Context, id int64) (*biz.User, error) {
	var user biz.User
	result := r.data.MySQL.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.log.Errorf("Failed to find user by id: %v", result.Error)
		return nil, result.Error
	}
	return &user, nil
}

// VerifyPassword 验证密码
func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
