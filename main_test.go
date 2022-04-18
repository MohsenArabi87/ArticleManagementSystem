package main_test

import (
	"MohsenArabi/ArticleManagementSystem/data"
	"MohsenArabi/ArticleManagementSystem/utils"
	"testing"
	"time"
)

func TestDeleteArticle(t *testing.T) {
	article := data.Article{Title: "sample title", Content: "This is a sample content", Tags: []string{"T1", "T2"}}
	logger := utils.NewLogger()
	repository := data.NewRepo(logger)
	newArticle, err := repository.CreateArticle(&article)
	if err != nil {
		t.Errorf(err.Error())
	} else {

		err1 := repository.DeleteArticle(article.ID)

		if err1 != nil {
			t.Errorf(err1.Error())
		} else {
			_, err2 := repository.GetArticleByID(newArticle.ID)
			if err2 == nil {
				t.Errorf(err2.Error())
			}
		}

	}
}
func TestUpdateArticle(t *testing.T) {
	article := data.Article{Title: "sample title", Content: "This is a sample content", Tags: []string{"T1", "T2"}}
	logger := utils.NewLogger()
	repository := data.NewRepo(logger)
	newArticle, err := repository.CreateArticle(&article)
	time.Sleep(1 * time.Second)
	if err != nil {
		t.Errorf(err.Error())
	} else {

		newArticle.Content = "Updated content"
		newArticle.Title = "updated title"
		_, err1 := repository.UpdateArticle(newArticle)

		if err1 != nil {
			t.Errorf(err1.Error())
		} else {
			fetchedArtecle, err2 := repository.GetArticleByID(newArticle.ID)
			if err2 != nil {
				t.Errorf(err2.Error())
			} else if !comapreArticles(fetchedArtecle, article) || !fetchedArtecle.UpdatedAt.After(fetchedArtecle.CreatedAt) {
				t.Errorf("Updated article does not match!")
				t.Log(fetchedArtecle.UpdatedAt)
				t.Log(fetchedArtecle.CreatedAt)
			}
		}

	}
}

func TestCreateArticle(t *testing.T) {
	article := data.Article{Title: "sample title", Content: "This is a sample content", Tags: []string{"T1", "T2"}}
	logger := utils.NewLogger()
	repository := data.NewRepo(logger)
	newArticle, err := repository.CreateArticle(&article)

	if err != nil {
		t.Errorf(err.Error())
	} else {

		fetchedArticle, err2 := repository.GetArticleByID(newArticle.ID)
		if err2 != nil {
			t.Errorf(err2.Error())
		} else if !comapreArticles(fetchedArticle, article) {
			t.Errorf("Inserted article does not match!")
		}

	}

}

func comapreArticles(article1 data.Article, article2 data.Article) bool {
	if article1.Author != article2.Author || article1.Title != article2.Title || article1.Content != article2.Content || article1.CreatedAt != article2.CreatedAt || len(article1.Tags) != len(article2.Tags) {
		return false
	}
	return true
}
