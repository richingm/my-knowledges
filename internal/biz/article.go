package biz

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/go-kratos/kratos/v2/log"
)

// Article 文章模型
type Article struct {
	ID              int64
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
	KnowledgeID     int64
	ParentArticleID int64
	Title           string
	Content         string
	Level           int32
	BySort          int64
}

// ArticleRepo 文章仓库接口
type ArticleRepo interface {
	Save(context.Context, *Article) (*Article, error)
	Update(context.Context, *Article) (*Article, error)
	FindByID(context.Context, int64) (*Article, error)
	Delete(context.Context, int64) error
	ListByKnowledgeID(context.Context, int64) ([]*Article, error)
}

// ArticleUsecase 文章用例
type ArticleUsecase struct {
	repo ArticleRepo
}

// NewArticleUsecase 创建文章用例
func NewArticleUsecase(repo ArticleRepo) *ArticleUsecase {
	return &ArticleUsecase{repo: repo}
}

// CreateArticle 创建文章
func (uc *ArticleUsecase) CreateArticle(ctx context.Context, article *Article) (*Article, error) {
	log.Infof("CreateArticle: %v", article.Title)
	now := time.Now()
	article.CreatedAt = now
	article.UpdatedAt = now

	article, err := uc.repo.Save(ctx, article)
	if err != nil {
		return nil, err
	}
	article.BySort = article.ID

	return uc.repo.Update(ctx, article)
}

// GetArticle 获取单个文章
func (uc *ArticleUsecase) GetArticle(ctx context.Context, id int64) (*Article, error) {
	return uc.repo.FindByID(ctx, id)
}

// UpdateArticle 更新文章
func (uc *ArticleUsecase) UpdateArticle(ctx context.Context, article *Article) (*Article, error) {
	log.Infof("UpdateArticle: %v", article.ID)
	// 查找文章
	oldArticle, err := uc.repo.FindByID(ctx, article.ID)
	if err != nil {
		return nil, err
	}

	// 更新文章信息
	if article.Title != "" {
		oldArticle.Title = article.Title
	}
	if article.Content != "" {
		oldArticle.Content = article.Content
	}
	if article.KnowledgeID != 0 {
		oldArticle.KnowledgeID = article.KnowledgeID
	}
	if article.ParentArticleID != 0 {
		oldArticle.ParentArticleID = article.ParentArticleID
	}
	if article.Level != 0 {
		oldArticle.Level = article.Level
	}
	if article.BySort != 0 {
		oldArticle.BySort = article.BySort
	}
	oldArticle.UpdatedAt = time.Now()
	return uc.repo.Update(ctx, oldArticle)
}

// DeleteArticle 删除文章
func (uc *ArticleUsecase) DeleteArticle(ctx context.Context, id int64) error {
	log.Infof("DeleteArticle: %v", id)
	return uc.repo.Delete(ctx, id)
}

// GetArticleTree 获取文章树
func (uc *ArticleUsecase) GetArticleTree(ctx context.Context, knowledgeID int64) ([]*Article, error) {
	return uc.repo.ListByKnowledgeID(ctx, knowledgeID)
}

// UpdateArticleSort 更新文章排序
func (uc *ArticleUsecase) UpdateArticleSort(ctx context.Context, id int64, bySort int64) error {
	// 获取文章
	article, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	// 更新排序值
	article.BySort = bySort
	article.UpdatedAt = time.Now()
	_, err = uc.repo.Update(ctx, article)
	return err
}

// MoveArticle 移动文章
func (uc *ArticleUsecase) MoveArticle(ctx context.Context, id int64, newParentId int64, knowledgeId int64) error {
	// 获取文章
	article, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	// 更新父文章 ID 和知识库 ID
	article.ParentArticleID = newParentId
	article.KnowledgeID = knowledgeId
	article.UpdatedAt = time.Now()
	_, err = uc.repo.Update(ctx, article)
	return err
}
