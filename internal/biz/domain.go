package biz

import (
	"context"
	"gorm.io/gorm"
	"time"
)

// Domain 领域模型
type Domain struct {
	ID          int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        string
	Description string
	BySort      int64
}

// DomainRepo 领域仓库接口
type DomainRepo interface {
	List(context.Context) ([]*Domain, error)
}

// DomainUsecase 领域用例
type DomainUsecase struct {
	repo DomainRepo
}

// NewDomainUsecase 创建领域用例
func NewDomainUsecase(repo DomainRepo) *DomainUsecase {
	return &DomainUsecase{repo: repo}
}

// ListDomains 获取领域列表
func (uc *DomainUsecase) ListDomains(ctx context.Context) ([]*Domain, error) {
	return uc.repo.List(ctx)
}
