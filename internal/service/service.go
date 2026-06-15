package service

import "github.com/google/wire"

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewGreeterService,
	NewDomainService,
	NewKnowledgeService,
	NewArticleService,
	NewFileService,
	NewUserService,
)
