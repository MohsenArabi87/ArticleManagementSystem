package service

import (
	"MohsenArabi/ArticleManagementSystem/utils"

	"github.com/hashicorp/go-hclog"
)

type Article interface {
}

// ArticleService is the implementation of our Article
type ArticleService struct {
	logger  hclog.Logger
	configs *utils.Configurations
}

// NewArticleService returns a new instance of the Article service
func NewArticleService(logger hclog.Logger, configs *utils.Configurations) *ArticleService {
	return &ArticleService{logger, configs}
}
