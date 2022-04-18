package handlers

import (
	"MohsenArabi/ArticleManagementSystem/data"
	"MohsenArabi/ArticleManagementSystem/service"
	"MohsenArabi/ArticleManagementSystem/utils"
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

// ArticleKey is used as a key for storing the Article object in context at middleware
type ArticleKey struct{}

// ArticleHandler wraps instances needed to perform operations on article object
type ArticleHandler struct {
	logger         hclog.Logger
	configs        *utils.Configurations
	validator      *data.Validation
	repo           data.Repository
	ArticleService service.Article
}

// NewArticleHandler returns a new ArticleHandler instance
func NewArticleHandler(l hclog.Logger, c *utils.Configurations, v *data.Validation, r data.Repository, articleSrvc service.Article) *ArticleHandler {
	return &ArticleHandler{
		logger:         l,
		configs:        c,
		validator:      v,
		repo:           r,
		ArticleService: articleSrvc,
	}
}

// MiddlewareValidateArticle validates the Article in the request
func (ah *ArticleHandler) MiddlewareValidateArticle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		ah.logger.Debug("Article json", r.Body)
		article := &data.Article{}

		err := data.FromJSON(article, r.Body)
		if err != nil {
			ah.logger.Error("deserialization of Article json failed", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericResponse{Status: false, Message: err.Error()}, w)
			return
		}

		// validate the user
		errs := ah.validator.Validate(article)
		if len(errs) != 0 {
			ah.logger.Error("validation of Article json failed", "error", errs)
			w.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericResponse{Status: false, Message: strings.Join(errs.Errors(), ",")}, w)
			return
		}

		// add the user to the context
		ctx := context.WithValue(r.Context(), ArticleKey{}, *article)
		r = r.WithContext(ctx)

		// call the next handler
		next.ServeHTTP(w, r)
	})
}

// CreateArticle handles CreateArticle request
func (ah *ArticleHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	article := r.Context().Value(ArticleKey{}).(data.Article)
	userID := r.Context().Value(UserIDKey{}).(string)
	article.Author = userID
	createdArticle, newErr := ah.repo.CreateArticle(&article)
	if newErr == nil {

		ah.logger.Debug("Article created successfully")
		w.WriteHeader(http.StatusCreated)
		data.ToJSON(&GenericResponse{Status: true, Message: "Article created successfully", Data: createdArticle}, w)
	}
}

//updateArticle handles updateArticle request
func (ah *ArticleHandler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	article := r.Context().Value(ArticleKey{}).(data.Article)
	userID := r.Context().Value(UserIDKey{}).(string)

	if article.Author == userID {
		updatedArticle, err := ah.repo.UpdateArticle(&article)
		if err == nil {

			ah.logger.Debug("Article updated successfully")
			w.WriteHeader(http.StatusCreated)
			data.ToJSON(&GenericResponse{Status: true, Message: "Article updated successfully", Data: updatedArticle}, w)
		}
	} else {
		ah.logger.Debug(utils.ErrCantUpdateOthersArticle)
		w.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericResponse{Status: false, Message: utils.ErrCantUpdateOthersArticle}, w)
	}
}

//DeleteArticle handles DeleteArticle request
func (ah *ArticleHandler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	articleID := params["articleID"]
	userID := r.Context().Value(UserIDKey{}).(string)
	article, err := ah.repo.GetArticleByID(articleID)

	if err != nil {
		ah.logger.Debug(utils.ErrArticleNotFound)
		w.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericResponse{Status: false, Message: utils.ErrArticleNotFound}, w)
	}

	if article.Author == userID {
		err := ah.repo.DeleteArticle(article.ID)
		if err == nil {

			ah.logger.Debug("Article Deleted successfully")
			w.WriteHeader(http.StatusCreated)
			data.ToJSON(&GenericResponse{Status: true, Message: "Article Deleted successfully"}, w)
		}
	} else {
		ah.logger.Debug(utils.ErrCantDeleteOthersArticle)
		w.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericResponse{Status: false, Message: utils.ErrCantDeleteOthersArticle}, w)
	}
}

//GetArticle handles getarticle request and fetch an article by id
func (ah *ArticleHandler) GetArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	articleID := params["articleID"]
	article, err := ah.repo.GetArticleByID(articleID)

	if err != nil {
		ah.logger.Debug(utils.ErrArticleNotFound)
		w.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericResponse{Status: false, Message: utils.ErrArticleNotFound}, w)
	} else {
		ah.logger.Debug("Article fetched successfully")
		w.WriteHeader(http.StatusCreated)
		data.ToJSON(&GenericResponse{Status: true, Message: "Article fetched successfully", Data: article}, w)
	}

}

//GetArticles handles GetArticles request and fetches a page of articles
func (ah *ArticleHandler) GetArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//params := mux.Vars(r)
	//pageNumber, err := strconv.Atoi(params["id"])
	params := r.FormValue("pageid")
	pageNumber, err := strconv.Atoi(params)

	if err != nil {
		ah.logger.Debug(utils.ErrInvalidPageNumber)
		w.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericResponse{Status: false, Message: utils.ErrInvalidPageNumber}, w)
		return
	}

	articles, err := ah.repo.GetArticles(pageNumber, ah.configs.PageSize)
	if err == nil {
		ah.logger.Debug("Article(s) fetched successfully")
		w.WriteHeader(http.StatusCreated)
		data.ToJSON(&GenericResponse{Status: true, Message: "Article(s) fetched successfully", Data: articles}, w)

	} else {
		ah.logger.Debug(utils.ErrInvalidPageNumber)
		w.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericResponse{Status: false, Message: utils.ErrInvalidPageNumber}, w)
	}
}

//GetArticlesTags handles GetArticlesTags requests and fetchs all existing tags
func (ah *ArticleHandler) GetArticlesTags(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tags := ah.repo.GetArticlesTags()
	ah.logger.Debug("Article tags fetched successfully")
	w.WriteHeader(http.StatusCreated)
	data.ToJSON(&GenericResponse{Status: true, Message: "Article tags fetched successfully", Data: tags}, w)

}
