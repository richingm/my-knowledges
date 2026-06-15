package service

import (
	"context"
	"errors"
	"time"

	v1 "my_knowledges/api/user/v1"
	"my_knowledges/internal/biz"
	"my_knowledges/internal/data"

	"github.com/go-kratos/kratos/v2/log"
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

// UserService 用户服务
type UserService struct {
	v1.UnimplementedUserServiceServer

	uc        *biz.UserUsecase
	log       *log.Helper
	jwtSecret string
}

// NewUserService 创建用户服务
func NewUserService(uc *biz.UserUsecase, logger log.Logger) *UserService {
	return &UserService{
		uc:        uc,
		log:       log.NewHelper(logger),
		jwtSecret: "your-secret-key-here", // 建议从配置中读取
	}
}

// SetJWTSecret 设置JWT密钥
func (s *UserService) SetJWTSecret(secret string) {
	s.jwtSecret = secret
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, in *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	user, err := s.uc.Register(ctx, in.Username, in.Email, in.Password)
	if err != nil {
		s.log.Errorf("Register failed: %v", err)
		return nil, err
	}

	token, err := s.generateToken(user.ID, user.Username)
	if err != nil {
		s.log.Errorf("Generate token failed: %v", err)
		return nil, err
	}

	return &v1.RegisterResponse{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Token:    token,
	}, nil
}

var ErrInvalidCredentials = errors.New("invalid username or password")

// Login 用户登录
func (s *UserService) Login(ctx context.Context, in *v1.LoginRequest) (*v1.LoginResponse, error) {
	user, err := s.uc.Login(ctx, in.Username, in.Password)
	if err != nil {
		s.log.Errorf("Login failed: %v", err)
		return nil, err
	}

	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if !data.VerifyPassword(user.Password, in.Password) {
		return nil, ErrInvalidCredentials
	}

	token, err := s.generateToken(user.ID, user.Username)
	if err != nil {
		s.log.Errorf("Generate token failed: %v", err)
		return nil, err
	}

	return &v1.LoginResponse{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Token:    token,
	}, nil
}

// GetCurrentUser 获取当前用户信息
func (s *UserService) GetCurrentUser(ctx context.Context, in *v1.GetCurrentUserRequest) (*v1.GetCurrentUserResponse, error) {
	return &v1.GetCurrentUserResponse{
		User: &v1.User{
			Id:       1,
			Username: "test",
			Email:    "test@example.com",
		},
	}, nil
}

// generateToken 生成JWT token
func (s *UserService) generateToken(uid int64, username string) (string, error) {
	claims := jwtv5.MapClaims{
		"uid":      uid,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // 7天有效期
	}

	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
