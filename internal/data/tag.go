package data

import (
	"context"

	"my_knowledges/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type tagRepo struct {
	data *Data
	log  *log.Helper
}

func NewTagRepo(data *Data, logger log.Logger) biz.TagRepo {
	return &tagRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *tagRepo) Save(ctx context.Context, tag *biz.Tag) (*biz.Tag, error) {
	r.log.Infof("Save tag: %s", tag.Name)
	result := r.data.MySQL.WithContext(ctx).Create(tag)
	if result.Error != nil {
		r.log.Errorf("Failed to save tag: %v", result.Error)
		return nil, result.Error
	}
	return tag, nil
}

func (r *tagRepo) Update(ctx context.Context, tag *biz.Tag) (*biz.Tag, error) {
	r.log.Infof("Update tag: %d", tag.ID)
	result := r.data.MySQL.WithContext(ctx).Model(tag).Updates(map[string]interface{}{
		"name":        tag.Name,
		"description": tag.Description,
		"updated_at":  tag.UpdatedAt,
	})
	if result.Error != nil {
		r.log.Errorf("Failed to update tag: %v", result.Error)
		return nil, result.Error
	}
	return tag, nil
}

func (r *tagRepo) Delete(ctx context.Context, id int64) error {
	r.log.Infof("Delete tag: %d", id)
	result := r.data.MySQL.WithContext(ctx).Model(&biz.Tag{}).Where("id = ?", id).Delete(&biz.Tag{})
	if result.Error != nil {
		r.log.Errorf("Failed to delete tag: %v", result.Error)
		return result.Error
	}
	return nil
}

func (r *tagRepo) FindByID(ctx context.Context, id int64) (*biz.Tag, error) {
	var tag biz.Tag
	result := r.data.MySQL.WithContext(ctx).Where("deleted_at IS NULL").Where("id = ?", id).First(&tag)
	if result.Error != nil {
		return nil, result.Error
	}
	return &tag, nil
}

func (r *tagRepo) FindByName(ctx context.Context, name string) (*biz.Tag, error) {
	var tag biz.Tag
	result := r.data.MySQL.WithContext(ctx).Where("deleted_at IS NULL").Where("name = ?", name).First(&tag)
	if result.Error != nil {
		return nil, result.Error
	}
	return &tag, nil
}

func (r *tagRepo) List(ctx context.Context) ([]*biz.Tag, error) {
	var tags []*biz.Tag
	result := r.data.MySQL.WithContext(ctx).Where("deleted_at IS NULL").Order("created_at DESC").Find(&tags)
	if result.Error != nil {
		return nil, result.Error
	}
	return tags, nil
}

func (r *tagRepo) SaveArticleTag(ctx context.Context, articleTag *biz.ArticleTag) error {
	r.log.Infof("Save article_tag: articleID=%d, tagID=%d", articleTag.ArticleID, articleTag.TagID)
	result := r.data.MySQL.WithContext(ctx).Create(articleTag)
	if result.Error != nil {
		r.log.Errorf("Failed to save article_tag: %v", result.Error)
		return result.Error
	}
	return nil
}

func (r *tagRepo) DeleteArticleTagsByArticleID(ctx context.Context, articleID int64) error {
	r.log.Infof("Delete article_tags by articleID: %d", articleID)
	result := r.data.MySQL.WithContext(ctx).Where("article_id = ?", articleID).Delete(&biz.ArticleTag{})
	if result.Error != nil {
		r.log.Errorf("Failed to delete article_tags: %v", result.Error)
		return result.Error
	}
	return nil
}

func (r *tagRepo) ListTagsByArticleID(ctx context.Context, articleID int64) ([]*biz.Tag, error) {
	var tags []*biz.Tag
	result := r.data.MySQL.WithContext(ctx).
		Joins("JOIN article_tags ON tags.id = article_tags.tag_id").
		Where("tags.deleted_at IS NULL").
		Where("article_tags.article_id = ?", articleID).
		Find(&tags)
	if result.Error != nil {
		return nil, result.Error
	}
	return tags, nil
}

func (r *tagRepo) ListArticlesByTagID(ctx context.Context, tagID int64) ([]*biz.Article, error) {
	var articles []*biz.Article
	result := r.data.MySQL.WithContext(ctx).
		Joins("JOIN article_tags ON articles.id = article_tags.article_id").
		Where("articles.deleted_at IS NULL").
		Where("article_tags.tag_id = ?", tagID).
		Order("articles.created_at DESC").
		Find(&articles)
	if result.Error != nil {
		return nil, result.Error
	}
	return articles, nil
}

func (r *tagRepo) UpdateTagArticleCount(ctx context.Context, tagID int64) error {
	r.log.Infof("Update tag article count: tagID=%d", tagID)
	var count int64
	result := r.data.MySQL.WithContext(ctx).
		Model(&biz.ArticleTag{}).
		Where("tag_id = ?", tagID).
		Count(&count)
	if result.Error != nil {
		r.log.Errorf("Failed to count article_tags: %v", result.Error)
		return result.Error
	}

	result = r.data.MySQL.WithContext(ctx).
		Model(&biz.Tag{}).
		Where("id = ?", tagID).
		Update("article_count", count)
	if result.Error != nil {
		r.log.Errorf("Failed to update tag article_count: %v", result.Error)
		return result.Error
	}
	return nil
}