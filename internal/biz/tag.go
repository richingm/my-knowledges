package biz

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/go-kratos/kratos/v2/log"
)

type Tag struct {
	ID           int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Name         string         `gorm:"unique;not null;size:100"`
	Description  string         `gorm:"size:500"`
	ArticleCount int64          `gorm:"default:0"`
}

type ArticleTag struct {
	ID        int64
	CreatedAt time.Time
	ArticleID int64 `gorm:"index"`
	TagID     int64 `gorm:"index"`
}

type TagRepo interface {
	Save(context.Context, *Tag) (*Tag, error)
	Update(context.Context, *Tag) (*Tag, error)
	Delete(context.Context, int64) error
	FindByID(context.Context, int64) (*Tag, error)
	FindByName(context.Context, string) (*Tag, error)
	List(context.Context) ([]*Tag, error)
	SaveArticleTag(context.Context, *ArticleTag) error
	DeleteArticleTagsByArticleID(context.Context, int64) error
	ListTagsByArticleID(context.Context, int64) ([]*Tag, error)
	ListArticlesByTagID(context.Context, int64) ([]*Article, error)
	UpdateTagArticleCount(context.Context, int64) error
}

type TagUsecase struct {
	repo TagRepo
}

func NewTagUsecase(repo TagRepo) *TagUsecase {
	return &TagUsecase{repo: repo}
}

func (uc *TagUsecase) CreateTag(ctx context.Context, tag *Tag) (*Tag, error) {
	log.Infof("CreateTag: %v", tag.Name)
	now := time.Now()
	tag.CreatedAt = now
	tag.UpdatedAt = now
	tag.ArticleCount = 0

	return uc.repo.Save(ctx, tag)
}

func (uc *TagUsecase) UpdateTag(ctx context.Context, tag *Tag) (*Tag, error) {
	log.Infof("UpdateTag: %v", tag.ID)
	oldTag, err := uc.repo.FindByID(ctx, tag.ID)
	if err != nil {
		return nil, err
	}

	if tag.Name != "" {
		oldTag.Name = tag.Name
	}
	if tag.Description != "" {
		oldTag.Description = tag.Description
	}
	oldTag.UpdatedAt = time.Now()

	return uc.repo.Update(ctx, oldTag)
}

func (uc *TagUsecase) DeleteTag(ctx context.Context, id int64) error {
	log.Infof("DeleteTag: %v", id)
	return uc.repo.Delete(ctx, id)
}

func (uc *TagUsecase) GetTag(ctx context.Context, id int64) (*Tag, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *TagUsecase) ListTags(ctx context.Context) ([]*Tag, error) {
	return uc.repo.List(ctx)
}

func (uc *TagUsecase) AddTagsToArticle(ctx context.Context, articleID int64, tagIDs []int64) error {
	log.Infof("AddTagsToArticle: articleID=%v, tagIDs=%v", articleID, tagIDs)

	if err := uc.repo.DeleteArticleTagsByArticleID(ctx, articleID); err != nil {
		return err
	}

	for _, tagID := range tagIDs {
		articleTag := &ArticleTag{
			ArticleID: articleID,
			TagID:     tagID,
		}
		if err := uc.repo.SaveArticleTag(ctx, articleTag); err != nil {
			return err
		}

		if err := uc.repo.UpdateTagArticleCount(ctx, tagID); err != nil {
			return err
		}
	}

	return nil
}

func (uc *TagUsecase) GetTagsByArticle(ctx context.Context, articleID int64) ([]*Tag, error) {
	return uc.repo.ListTagsByArticleID(ctx, articleID)
}

func (uc *TagUsecase) GetArticlesByTag(ctx context.Context, tagID int64) ([]*Article, error) {
	return uc.repo.ListArticlesByTagID(ctx, tagID)
}
