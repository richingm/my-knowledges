package server

import (
	filev1 "my_knowledges/api/file/v1"
	v1 "my_knowledges/api/helloworld/v1"
	knowledgev1 "my_knowledges/api/knowledge/v1"
	userv1 "my_knowledges/api/user/v1"
	"my_knowledges/internal/conf"
	"my_knowledges/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, greeter *service.GreeterService, domain *service.DomainService, knowledge *service.KnowledgeService, article *service.ArticleService, file *service.FileService, user *service.UserService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterGreeterServer(srv, greeter)
	knowledgev1.RegisterDomainServiceServer(srv, domain)
	knowledgev1.RegisterKnowledgeServiceServer(srv, knowledge)
	knowledgev1.RegisterArticleServiceServer(srv, article)
	filev1.RegisterFileServiceServer(srv, file)
	userv1.RegisterUserServiceServer(srv, user)
	return srv
}
