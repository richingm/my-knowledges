package data

import (
	"context"
	"time"

	"my_knowledges/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type knowledgeRepo struct {
	data *Data
	log  *log.Helper
}

// NewKnowledgeRepo 创建知识库仓库
func NewKnowledgeRepo(data *Data, logger log.Logger) biz.KnowledgeRepo {
	return &knowledgeRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *knowledgeRepo) Save(ctx context.Context, knowledge *biz.Knowledge) (*biz.Knowledge, error) {
	r.log.Infof("Save knowledge base: %s", knowledge.Name)
	result := r.data.MySQL.WithContext(ctx).Create(knowledge)
	if result.Error != nil {
		r.log.Errorf("Failed to save knowledge: %v", result.Error)
		return nil, result.Error
	}
	return knowledge, nil
}

func (r *knowledgeRepo) Update(ctx context.Context, knowledge *biz.Knowledge) (*biz.Knowledge, error) {
	r.log.Infof("Update knowledge base: %d", knowledge.ID)
	// 只更新名称、描述等字段，不修改 parent_knowledge_id（通过 move 接口修改）
	result := r.data.MySQL.WithContext(ctx).Model(knowledge).Updates(map[string]interface{}{
		"domain_id":   knowledge.DomainID,
		"name":        knowledge.Name,
		"description": knowledge.Description,
		"by_sort":     knowledge.BySort,
		"updated_at":  knowledge.UpdatedAt,
	})
	if result.Error != nil {
		r.log.Errorf("Failed to update knowledge: %v1", result.Error)
		return nil, result.Error
	}
	return knowledge, nil
}

func (r *knowledgeRepo) Delete(ctx context.Context, id int64) error {
	result := r.data.MySQL.WithContext(ctx).Model(&biz.Knowledge{}).Where("id = ?", id).Delete(&biz.Knowledge{})
	if result.Error != nil {
		r.log.Errorf("Failed to delete knowledge base: %v1", result.Error)
		return result.Error
	}
	return nil
}

func (r *knowledgeRepo) ListByDomainID(ctx context.Context, domainID int64) ([]*biz.Knowledge, error) {
	var knowledges []*biz.Knowledge
	result := r.data.MySQL.WithContext(ctx).Where("deleted_at IS NULL").Where("domain_id = ?", domainID).Order("by_sort ASC").Find(&knowledges)
	if result.Error != nil {
		return nil, result.Error
	}
	return knowledges, nil
}

func (r *knowledgeRepo) Get(ctx context.Context, id int64) (*biz.Knowledge, error) {
	var knowledge biz.Knowledge
	result := r.data.MySQL.WithContext(ctx).Where("id = ?", id).First(&knowledge)
	if result.Error != nil {
		r.log.Errorf("Failed to get knowledge: %v", result.Error)
		return nil, result.Error
	}
	return &knowledge, nil
}

func (r *knowledgeRepo) UpdateParentAndSort(ctx context.Context, id int64, parentID int64, bySort int64) error {
	result := r.data.MySQL.WithContext(ctx).Model(&biz.Knowledge{}).Where("id = ?", id).Updates(map[string]interface{}{
		"parent_knowledge_id": parentID,
		"by_sort":             bySort,
		"updated_at":          time.Now(),
	})
	if result.Error != nil {
		r.log.Errorf("Failed to update parent and sort: %v", result.Error)
		return result.Error
	}
	return nil
}
