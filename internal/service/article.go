package service

import (
	"context"

	v1 "my_knowledges/api/knowledge/v1"
	"my_knowledges/internal/biz"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// ArticleService 文章服务
type ArticleService struct {
	v1.UnimplementedArticleServiceServer

	uc *biz.ArticleUsecase
}

// NewArticleService 创建文章服务
func NewArticleService(uc *biz.ArticleUsecase) *ArticleService {
	return &ArticleService{uc: uc}
}

// CreateArticle 创建文章
func (s *ArticleService) CreateArticle(ctx context.Context, in *v1.CreateArticleRequest) (*v1.ArticleResponse, error) {
	article, err := s.uc.CreateArticle(ctx, &biz.Article{
		KnowledgeID:     in.KnowledgeId,
		ParentArticleID: in.ParentArticleId,
		Title:           in.Title,
		Content:         in.Content,
		Level:           in.Level,
		BySort:          in.BySort,
	})
	if err != nil {
		return nil, err
	}

	return &v1.ArticleResponse{
		Id:              article.ID,
		CreatedAt:       timestamppb.New(article.CreatedAt),
		UpdatedAt:       timestamppb.New(article.UpdatedAt),
		KnowledgeId:     article.KnowledgeID,
		ParentArticleId: article.ParentArticleID,
		Title:           article.Title,
		Content:         article.Content,
		Level:           article.Level,
		BySort:          article.BySort,
	}, nil
}

// GetArticle 获取单个文章
func (s *ArticleService) GetArticle(ctx context.Context, in *v1.GetArticleRequest) (*v1.Article, error) {
	article, err := s.uc.GetArticle(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	var deletedAt *timestamppb.Timestamp
	if article.DeletedAt.Valid {
		deletedAt = timestamppb.New(article.DeletedAt.Time)
	}
	return &v1.Article{
		Id:              article.ID,
		CreatedAt:       timestamppb.New(article.CreatedAt),
		UpdatedAt:       timestamppb.New(article.UpdatedAt),
		DeletedAt:       deletedAt,
		KnowledgeId:     article.KnowledgeID,
		ParentArticleId: article.ParentArticleID,
		Title:           article.Title,
		Content:         article.Content,
		Level:           article.Level,
		BySort:          article.BySort,
	}, nil
}

// UpdateArticle 更新文章
func (s *ArticleService) UpdateArticle(ctx context.Context, in *v1.UpdateArticleRequest) (*v1.ArticleResponse, error) {
	article, err := s.uc.UpdateArticle(ctx, &biz.Article{
		ID:              in.Id,
		KnowledgeID:     in.KnowledgeId,
		ParentArticleID: in.ParentArticleId,
		Title:           in.Title,
		Content:         in.Content,
		Level:           in.Level,
		BySort:          in.BySort,
	})
	if err != nil {
		return nil, err
	}

	return &v1.ArticleResponse{
		Id:              article.ID,
		CreatedAt:       timestamppb.New(article.CreatedAt),
		UpdatedAt:       timestamppb.New(article.UpdatedAt),
		KnowledgeId:     article.KnowledgeID,
		ParentArticleId: article.ParentArticleID,
		Title:           article.Title,
		Content:         article.Content,
		Level:           article.Level,
		BySort:          article.BySort,
	}, nil
}

// DeleteArticle 删除文章
func (s *ArticleService) DeleteArticle(ctx context.Context, in *v1.DeleteArticleRequest) (*v1.ArticleResponse, error) {
	// 删除文章
	err := s.uc.DeleteArticle(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &v1.ArticleResponse{}, nil
}

// GetArticleTree 获取文章树
func (s *ArticleService) GetArticleTree(ctx context.Context, in *v1.GetArticleTreeRequest) (*v1.GetArticleTreeResponse, error) {
	articles, err := s.uc.GetArticleTree(ctx, in.KnowledgeId)
	if err != nil {
		return nil, err
	}

	// 构建文章树
	nodeMap := make(map[int64]*v1.ArticleTreeNode)
	var rootNodes []*v1.ArticleTreeNode

	// 首先创建所有节点
	for _, article := range articles {
		node := &v1.ArticleTreeNode{
			Id:              article.ID,
			KnowledgeId:     article.KnowledgeID,
			ParentArticleId: article.ParentArticleID,
			Title:           article.Title,
			Content:         article.Content,
			Level:           article.Level,
			BySort:          article.BySort,
			Children:        []*v1.ArticleTreeNode{},
		}
		nodeMap[article.ID] = node

		// 根节点（parentArticleId 为 0）
		if article.ParentArticleID == 0 {
			rootNodes = append(rootNodes, node)
		}
	}

	// 构建父子关系
	for _, article := range articles {
		if article.ParentArticleID != 0 {
			if parentNode, ok := nodeMap[article.ParentArticleID]; ok {
				parentNode.Children = append(parentNode.Children, nodeMap[article.ID])
			}
		}
	}

	return &v1.GetArticleTreeResponse{
		Items: rootNodes,
	}, nil
}

// SortArticle 排序文章
func (s *ArticleService) SortArticle(ctx context.Context, in *v1.SortArticleRequest) (*v1.SortArticleResponse, error) {
	// 遍历排序项，更新文章排序值
	var knowledgeID int64
	for _, item := range in.Items {
		// 获取文章信息，以获取知识库 ID
		article, err := s.uc.GetArticle(ctx, item.Id)
		if err != nil {
			return &v1.SortArticleResponse{
				Success: false,
				Message: "排序失败: " + err.Error(),
			}, err
		}
		knowledgeID = article.KnowledgeID

		// 更新排序值
		err = s.uc.UpdateArticleSort(ctx, item.Id, item.BySort)
		if err != nil {
			return &v1.SortArticleResponse{
				Success: false,
				Message: "排序失败: " + err.Error(),
			}, err
		}
	}

	// 获取文章树
	articleTree, err := s.GetArticleTree(ctx, &v1.GetArticleTreeRequest{
		KnowledgeId: knowledgeID,
	})
	if err != nil {
		return &v1.SortArticleResponse{
			Success: false,
			Message: "获取文章树失败: " + err.Error(),
		}, err
	}

	return &v1.SortArticleResponse{
		Success:     true,
		Message:     "排序成功",
		ArticleTree: articleTree.Items,
	}, nil
}

// MoveArticle 移动文章
func (s *ArticleService) MoveArticle(ctx context.Context, in *v1.MoveArticleRequest) (*v1.MoveArticleResponse, error) {
	// 移动文章
	err := s.uc.MoveArticle(ctx, in.Id, in.NewParentId, in.KnowledgeId)
	if err != nil {
		return &v1.MoveArticleResponse{
			Success: false,
			Message: "移动失败: " + err.Error(),
		}, err
	}

	// 获取文章树
	articleTree, err := s.GetArticleTree(ctx, &v1.GetArticleTreeRequest{
		KnowledgeId: in.KnowledgeId,
	})
	if err != nil {
		return &v1.MoveArticleResponse{
			Success: false,
			Message: "获取文章树失败: " + err.Error(),
		}, err
	}

	return &v1.MoveArticleResponse{
		Success:     true,
		Message:     "移动成功",
		ArticleTree: articleTree.Items,
	}, nil
}
