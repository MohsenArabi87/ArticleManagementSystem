package data

// Repository is an interface for the storage implementation of services
type Repository interface {
	Create(user *User) error
	GetUserByEmail(email string) (*User, error)
	CreateArticle(article *Article) (*Article, error)
	UpdateArticle(article *Article) (*Article, error)
	DeleteArticle(articleID string) error
	GetArticles(pageNumber int, pagesize int) ([]Article, error)
	GetArticleByID(articleID string) (Article, error)
	GetArticlesTags() []string
}
