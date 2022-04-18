package data

import (
	"MohsenArabi/ArticleManagementSystem/utils"
	"errors"
	"time"

	"github.com/hashicorp/go-hclog"
	uuid "github.com/satori/go.uuid"
)

var users = make(map[string]User)
var articles = make(map[string]Article)
var tags []string

// Repo has the implementation of the in memory repository.
type Repo struct {
	logger hclog.Logger
}

// NewRepo returns a new Repo instance
func NewRepo(logger hclog.Logger) *Repo {
	return &Repo{logger}
}

//creates a new user
func (repo *Repo) Create(user *User) error {
	if _, exists := users[user.Email]; exists {
		repo.logger.Info(utils.ErrUserAlreadyExists)
		return errors.New(utils.ErrUserAlreadyExists)
	} else {
		repo.logger.Info("creating user", hclog.Fmt("%#v", user))
		users[user.Email] = *user
		return nil
	}
}

//get user struct by email
func (repo *Repo) GetUserByEmail(email string) (*User, error) {
	repo.logger.Debug("searching for user with email", email)

	if u, exists := users[email]; exists {
		repo.logger.Debug("read users", hclog.Fmt("%#v", u))
		return &u, nil
	} else {
		return nil, errors.New(utils.ErrUserNotFound)
	}
}

// creates new article
func (repo *Repo) CreateArticle(article *Article) (*Article, error) {
	repo.logger.Info("creating article")
	article.ID = uuid.NewV4().String()
	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()
	articles[article.ID] = *article
	return article, nil
}

//updates only title, content and tags of article
func (repo *Repo) UpdateArticle(newArticle *Article) (*Article, error) {
	repo.logger.Info("updating article")
	if oldArticle, exists := articles[newArticle.ID]; exists {
		oldArticle.UpdatedAt = time.Now()
		oldArticle.Title = newArticle.Title
		oldArticle.Content = newArticle.Content
		oldArticle.Tags = newArticle.Tags
		articles[newArticle.ID] = oldArticle
		return &oldArticle, nil

	} else {
		return nil, errors.New(utils.ErrArticleNotFound)
	}

}

//deletes the article by ID
func (repo *Repo) DeleteArticle(articleID string) error {
	repo.logger.Info("deleting article")
	if _, exists := articles[articleID]; exists {
		delete(articles, articleID)
		return nil

	} else {
		return errors.New(utils.ErrArticleNotFound)
	}

}

//fetchs only one page of articles
func (repo *Repo) GetArticles(pageNumber int, pageSize int) ([]Article, error) {
	repo.logger.Info(("fetching articles"))
	start := (pageNumber - 1) * pageSize
	stop := start + pageSize

	if start >= len(articles) {
		return nil, errors.New(utils.ErrInvalidPageNumber)
	}

	if stop > len(articles) {
		stop = len(articles)
	}
	repo.logger.Debug("fetching articles from %v to %v", start, stop)
	i := 0
	var result []Article
	for _, article := range articles {
		if i >= start && i < stop {
			result = append(result, article)
		}
		i++
		if i >= stop {
			break
		}
	}
	return result, nil

}

//get an article by ID
func (repo *Repo) GetArticleByID(articleID string) (Article, error) {
	repo.logger.Info(("fetching article"))
	article, exists := articles[articleID]
	if exists {
		return article, nil
	} else {
		return article, errors.New(utils.ErrArticleNotFound)
	}
}

//get all unique article tags
func (repo *Repo) GetArticlesTags() []string {
	repo.logger.Info(("fetching article tags"))
	tagsMap := make(map[string]struct{})
	tags := []string{}
	var empty struct{}
	for _, article := range articles {
		for _, tag := range article.Tags {
			if _, exists := tagsMap[tag]; !exists {
				tagsMap[tag] = empty
				tags = append(tags, tag)
			}
		}
	}
	return tags

}
