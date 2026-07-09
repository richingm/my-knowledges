package biz

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/go-kratos/kratos/v2/log"
)

// Knowledge 知识库模型
type Knowledge struct {
	ID                int64
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
	DomainID          int64
	ParentKnowledgeID int64
	Name              string
	Description       string
	BySort            int64
}

// KnowledgeRepo 知识库仓库接口
type KnowledgeRepo interface {
	Save(context.Context, *Knowledge) (*Knowledge, error)
	Update(context.Context, *Knowledge) (*Knowledge, error)
	UpdateParentAndSort(context.Context, int64, int64, int64) error
	Delete(context.Context, int64) error
	Get(context.Context, int64) (*Knowledge, error)
	ListByDomainID(context.Context, int64) ([]*Knowledge, error)
}

// KnowledgeUsecase 知识库用例
type KnowledgeUsecase struct {
	repo KnowledgeRepo
}

// NewKnowledgeUsecase 创建知识库用例
func NewKnowledgeUsecase(repo KnowledgeRepo) *KnowledgeUsecase {
	return &KnowledgeUsecase{repo: repo}
}

// CreateKnowledge 创建知识库
func (uc *KnowledgeUsecase) CreateKnowledge(ctx context.Context, knowledge *Knowledge) (*Knowledge, error) {
	log.Infof("CreateKnowledge: %v", knowledge.Name)
	log.Infof("CreateKnowledge: %v", knowledge.ParentKnowledgeID)
	now := time.Now()
	knowledge.CreatedAt = now
	knowledge.UpdatedAt = now
	return uc.repo.Save(ctx, knowledge)
}

// UpdateKnowledge 更新知识库
func (uc *KnowledgeUsecase) UpdateKnowledge(ctx context.Context, knowledge *Knowledge) (*Knowledge, error) {
	log.Infof("UpdateKnowledge: %v", knowledge.ID)
	knowledge.UpdatedAt = time.Now()
	return uc.repo.Update(ctx, knowledge)
}

// DeleteKnowledge 删除知识库
func (uc *KnowledgeUsecase) DeleteKnowledge(ctx context.Context, id int64) error {
	log.Infof("DeleteKnowledge: %v", id)
	return uc.repo.Delete(ctx, id)
}

// GetKnowledgeTree 获取知识库树
func (uc *KnowledgeUsecase) GetKnowledgeTree(ctx context.Context, domainID int64) ([]*Knowledge, error) {
	return uc.repo.ListByDomainID(ctx, domainID)
}

// UpdateKnowledgeSort 更新知识库排序
func (uc *KnowledgeUsecase) UpdateKnowledgeSort(ctx context.Context, id int64, bySort int64) error {
	// 获取知识库
	knowledge, err := uc.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	// 更新排序值
	knowledge.BySort = bySort
	knowledge.UpdatedAt = time.Now()
	_, err = uc.repo.Update(ctx, knowledge)
	return err
}

// MoveKnowledge 移动知识库
func (uc *KnowledgeUsecase) MoveKnowledge(ctx context.Context, id int64, newParentId int64, domainId int64) error {
	// 获取知识库
	knowledge, err := uc.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	// 不能移动到自身或自身的子节点下
	if id == newParentId {
		return fmt.Errorf("不能移动到自身")
	}
	// 更新父知识库 ID 和排序
	knowledge.ParentKnowledgeID = newParentId
	knowledge.DomainID = domainId
	knowledge.UpdatedAt = time.Now()
	return uc.repo.UpdateParentAndSort(ctx, id, newParentId, knowledge.BySort)
}

// MoveKnowledgeWithSort 移动知识库并设置排序
func (uc *KnowledgeUsecase) MoveKnowledgeWithSort(ctx context.Context, id int64, newParentId int64, bySort int64) error {
	if id == newParentId {
		return fmt.Errorf("不能移动到自身")
	}
	return uc.repo.UpdateParentAndSort(ctx, id, newParentId, bySort)
}
