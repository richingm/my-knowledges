package data

import (
	"context"

	"my_knowledges/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type domainRepo struct {
	data *Data
	log  *log.Helper
}

// NewDomainRepo 创建领域仓库
func NewDomainRepo(data *Data, logger log.Logger) biz.DomainRepo {
	return &domainRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *domainRepo) List(ctx context.Context) ([]*biz.Domain, error) {
	var domains []*biz.Domain
	result := r.data.MySQL.WithContext(ctx).Where("deleted_at IS NULL").Order("by_sort ASC").Find(&domains)
	if result.Error != nil {
		r.log.Errorf("Failed to list knowledgebase groups: %v1", result.Error)
		return nil, result.Error
	}
	return domains, nil
}
