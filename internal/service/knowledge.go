package service

import (
	"context"

	v1 "my_knowledges/api/knowledge/v1"
	"my_knowledges/internal/biz"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// KnowledgeService 知识库服务
type KnowledgeService struct {
	v1.UnimplementedKnowledgeServiceServer

	uc *biz.KnowledgeUsecase
}

// NewKnowledgeService 创建知识库服务
func NewKnowledgeService(uc *biz.KnowledgeUsecase) *KnowledgeService {
	return &KnowledgeService{uc: uc}
}

// CreateKnowledge 创建知识库
func (s *KnowledgeService) CreateKnowledge(ctx context.Context, in *v1.CreateKnowledgeRequest) (*v1.Knowledge, error) {
	knowledge, err := s.uc.CreateKnowledge(ctx, &biz.Knowledge{
		DomainID:          in.DomainId,
		ParentKnowledgeID: in.ParentKnowledgeId,
		Name:              in.Name,
		Description:       in.Description,
		BySort:            in.BySort,
	})
	if err != nil {
		return nil, err
	}
	var deletedAt *timestamppb.Timestamp
	if knowledge.DeletedAt.Valid {
		deletedAt = timestamppb.New(knowledge.DeletedAt.Time)
	}
	return &v1.Knowledge{
		Id:                knowledge.ID,
		CreatedAt:         timestamppb.New(knowledge.CreatedAt),
		UpdatedAt:         timestamppb.New(knowledge.UpdatedAt),
		DeletedAt:         deletedAt,
		DomainId:          knowledge.DomainID,
		ParentKnowledgeId: knowledge.ParentKnowledgeID,
		Name:              knowledge.Name,
		Description:       knowledge.Description,
		BySort:            knowledge.BySort,
	}, nil
}

// UpdateKnowledge 更新知识库
func (s *KnowledgeService) UpdateKnowledge(ctx context.Context, in *v1.UpdateKnowledgeRequest) (*v1.Knowledge, error) {
	knowledge, err := s.uc.UpdateKnowledge(ctx, &biz.Knowledge{
		ID:                in.Id,
		DomainID:          in.DomainId,
		ParentKnowledgeID: in.ParentKnowledgeId,
		Name:              in.Name,
		Description:       in.Description,
		BySort:            in.BySort,
	})
	if err != nil {
		return nil, err
	}
	var deletedAt *timestamppb.Timestamp
	if knowledge.DeletedAt.Valid {
		deletedAt = timestamppb.New(knowledge.DeletedAt.Time)
	}
	return &v1.Knowledge{
		Id:                knowledge.ID,
		CreatedAt:         timestamppb.New(knowledge.CreatedAt),
		UpdatedAt:         timestamppb.New(knowledge.UpdatedAt),
		DeletedAt:         deletedAt,
		DomainId:          knowledge.DomainID,
		ParentKnowledgeId: knowledge.ParentKnowledgeID,
		Name:              knowledge.Name,
		Description:       knowledge.Description,
		BySort:            knowledge.BySort,
	}, nil
}

// DeleteKnowledge 删除知识库
func (s *KnowledgeService) DeleteKnowledge(ctx context.Context, in *v1.DeleteKnowledgeRequest) (*emptypb.Empty, error) {
	err := s.uc.DeleteKnowledge(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// GetKnowledgeTree 获取知识库树
func (s *KnowledgeService) GetKnowledgeTree(ctx context.Context, in *v1.GetKnowledgeTreeRequest) (*v1.GetKnowledgeTreeResponse, error) {
	knowledges, err := s.uc.GetKnowledgeTree(ctx, in.DomainId)
	if err != nil {
		return nil, err
	}

	// 构建知识库树
	nodeMap := make(map[int64]*v1.KnowledgeTreeNode)
	var rootNodes []*v1.KnowledgeTreeNode

	// 首先创建所有节点
	for _, knowledge := range knowledges {
		node := &v1.KnowledgeTreeNode{
			Id:                knowledge.ID,
			DomainId:          knowledge.DomainID,
			ParentKnowledgeId: knowledge.ParentKnowledgeID,
			Name:              knowledge.Name,
			Description:       knowledge.Description,
			BySort:            knowledge.BySort,
			Children:          []*v1.KnowledgeTreeNode{},
		}
		nodeMap[knowledge.ID] = node

		// 根节点（parentKnowledgeId 为 0）
		if knowledge.ParentKnowledgeID == 0 {
			rootNodes = append(rootNodes, node)
		}
	}

	// 构建父子关系
	for _, knowledge := range knowledges {
		if knowledge.ParentKnowledgeID != 0 {
			if parentNode, ok := nodeMap[knowledge.ParentKnowledgeID]; ok {
				parentNode.Children = append(parentNode.Children, nodeMap[knowledge.ID])
			}
		}
	}

	return &v1.GetKnowledgeTreeResponse{
		Items: rootNodes,
	}, nil
}

// SortKnowledge 排序知识库
func (s *KnowledgeService) SortKnowledge(ctx context.Context, in *v1.SortKnowledgeRequest) (*v1.SortKnowledgeResponse, error) {
	// 遍历排序项，更新知识库排序值
	for _, item := range in.Items {
		err := s.uc.UpdateKnowledgeSort(ctx, item.Id, item.BySort)
		if err != nil {
			return &v1.SortKnowledgeResponse{
				Success: false,
				Message: "排序失败: " + err.Error(),
			}, err
		}
	}

	return &v1.SortKnowledgeResponse{
		Success: true,
		Message: "排序成功",
	}, nil
}

// MoveKnowledge 移动知识库
func (s *KnowledgeService) MoveKnowledge(ctx context.Context, in *v1.MoveKnowledgeRequest) (*v1.MoveKnowledgeResponse, error) {
	// 移动知识库
	err := s.uc.MoveKnowledge(ctx, in.Id, in.NewParentId, in.DomainId)
	if err != nil {
		return &v1.MoveKnowledgeResponse{
			Success: false,
			Message: "移动失败: " + err.Error(),
		}, err
	}

	// 获取知识库树
	knowledgeTree, err := s.GetKnowledgeTree(ctx, &v1.GetKnowledgeTreeRequest{
		DomainId: in.DomainId,
	})
	if err != nil {
		return &v1.MoveKnowledgeResponse{
			Success: false,
			Message: "获取知识库树失败: " + err.Error(),
		}, err
	}

	return &v1.MoveKnowledgeResponse{
		Success:      true,
		Message:      "移动成功",
		KnowledgeTree: knowledgeTree.Items,
	}, nil
}
