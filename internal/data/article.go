package data

import (
	"context"

	"my_knowledges/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type articleRepo struct {
	data *Data
	log  *log.Helper
}

// NewArticleRepo 创建文章仓库
func NewArticleRepo(data *Data, logger log.Logger) biz.ArticleRepo {
	return &articleRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *articleRepo) Save(ctx context.Context, article *biz.Article) (*biz.Article, error) {
	r.log.Infof("Save article: %s", article.Title)
	result := r.data.MySQL.WithContext(ctx).Create(article)
	if result.Error != nil {
		r.log.Errorf("Failed to save article: %v1", result.Error)
		return nil, result.Error
	}
	return article, nil
}

func (r *articleRepo) Update(ctx context.Context, article *biz.Article) (*biz.Article, error) {
	r.log.Infof("Update article: %d", article.ID)
	// 只更新需要的字段，避免修改 CreatedAt 字段
	result := r.data.MySQL.WithContext(ctx).Model(article).Updates(map[string]interface{}{
		"knowledge_id":      article.KnowledgeID,
		"parent_article_id": article.ParentArticleID,
		"title":             article.Title,
		"content":           article.Content,
		"level":             article.Level,
		"by_sort":           article.BySort,
		"updated_at":        article.UpdatedAt,
	})
	if result.Error != nil {
		r.log.Errorf("Failed to update article: %v1", result.Error)
		return nil, result.Error
	}
	return article, nil
}

func (r *articleRepo) FindByID(ctx context.Context, id int64) (*biz.Article, error) {
	var article *biz.Article
	result := r.data.MySQL.WithContext(ctx).Where("deleted_at IS NULL").Where("id = ?", id).Find(&article)
	if result.Error != nil {
		return nil, result.Error
	}
	return article, nil
}

func (r *articleRepo) Delete(ctx context.Context, id int64) error {
	r.log.Infof("Delete article: %d", id)
	result := r.data.MySQL.WithContext(ctx).Model(&biz.Article{}).Where("id = ?", id).Delete(&biz.Article{})
	if result.Error != nil {
		r.log.Errorf("Failed to delete article: %v1", result.Error)
		return result.Error
	}
	return nil
}

func (r *articleRepo) ListByKnowledgeID(ctx context.Context, knowledgeID int64) ([]*biz.Article, error) {
	var articles []*biz.Article
	result := r.data.MySQL.WithContext(ctx).Where("deleted_at IS NULL").Where("knowledge_id = ?", knowledgeID).Order("by_sort ASC").Find(&articles)
	if result.Error != nil {
		return nil, result.Error
	}
	return articles, nil
}
