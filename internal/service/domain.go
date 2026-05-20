package service

import (
	"context"

	v1 "my_knowledges/api/knowledge/v1"
	"my_knowledges/internal/biz"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// DomainService 领域服务
type DomainService struct {
	v1.UnimplementedDomainServiceServer

	uc *biz.DomainUsecase
}

// NewDomainService 创建领域服务
func NewDomainService(uc *biz.DomainUsecase) *DomainService {
	return &DomainService{uc: uc}
}

// ListDomains 获取领域列表
func (s *DomainService) ListDomains(ctx context.Context, in *v1.ListDomainsRequest) (*v1.ListDomainsResponse, error) {
	domains, err := s.uc.ListDomains(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]*v1.Domain, len(domains))
	for i, domain := range domains {
		var deletedAt *timestamppb.Timestamp
		if domain.DeletedAt.Valid {
			deletedAt = timestamppb.New(domain.DeletedAt.Time)
		}
		items[i] = &v1.Domain{
			Id:          domain.ID,
			CreatedAt:   timestamppb.New(domain.CreatedAt),
			UpdatedAt:   timestamppb.New(domain.UpdatedAt),
			DeletedAt:   deletedAt,
			Name:        domain.Name,
			Description: domain.Description,
			BySort:      domain.BySort,
		}
	}
	return &v1.ListDomainsResponse{
		Items: items,
	}, nil
}
