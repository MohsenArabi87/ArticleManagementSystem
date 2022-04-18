package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"MohsenArabi/ArticleManagementSystem/data"
	"MohsenArabi/ArticleManagementSystem/handlers"
	"MohsenArabi/ArticleManagementSystem/service"
	"MohsenArabi/ArticleManagementSystem/utils"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

func main() {

	//creates new insatnce of Logger
	logger := utils.NewLogger()

	//creates new insatnce of configurations
	configs := utils.NewConfigurations(logger)

	// validator contains all the methods that are need to validate the user json in request
	validator := data.NewValidation()

	// repository contains all the methods that interact with structs to perform CURD operations.
	repository := data.NewRepo(logger)

	// authService contains all methods that help in authorizing a user request
	authService := service.NewAuthService(logger, configs)

	// articleService contains all methods that help in managing articles
	articleService := service.NewArticleService(logger, configs)

	// UserHandler encapsulates all the requests related to user
	uh := handlers.NewAuthHandler(logger, configs, validator, repository, authService)
	// ArticleHandler encapsulates all the requests related to article
	ah := handlers.NewArticleHandler(logger, configs, validator, repository, articleService)

	// create a serve mux
	sm := mux.NewRouter()

	// register handlers
	postR := sm.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/signup", uh.Signup)
	postR.HandleFunc("/login", uh.Login)
	postR.Use(uh.MiddlewareValidateUser)

	// used the PathPrefix as workaround for scenarios where all the
	// get requests must use the ValidateAccessToken middleware except
	// the /refresh-token request which has to use ValidateRefreshToken middleware
	refToken := sm.PathPrefix("/refresh-token").Subrouter()
	refToken.HandleFunc("", uh.RefreshToken)
	refToken.Use(uh.MiddlewareValidateRefreshToken)

	//handlers for creating and updating an article and validates article and access token at middleware
	postRArticles := sm.PathPrefix("/Article").Methods(http.MethodPost).Subrouter()
	postRArticles.HandleFunc("/Create", ah.CreateArticle)
	postRArticles.HandleFunc("/Update", ah.UpdateArticle)
	postRArticles.Use(uh.MiddlewareValidateAccessToken)
	postRArticles.Use(ah.MiddlewareValidateArticle)

	//handlers for fetching and deleting article and validates access token at middleware
	getArticles := sm.PathPrefix("/Article").Methods(http.MethodGet).Subrouter()
	getArticles.HandleFunc("/Tags", ah.GetArticlesTags)
	getArticles.HandleFunc("/Delete/{articleID}", ah.DeleteArticle)
	getArticles.HandleFunc("", ah.GetArticles).Queries("pageid", "{id:[0-9]+}")
	getArticles.HandleFunc("/{articleID}", ah.GetArticle)
	getArticles.Use(uh.MiddlewareValidateAccessToken)

	// create a server
	svr := http.Server{
		Addr:         configs.ServerAddress,
		Handler:      sm,
		ErrorLog:     logger.StandardLogger(&hclog.StandardLoggerOptions{}),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// start the server
	go func() {
		logger.Info("starting the server at port", configs.ServerAddress)

		err := svr.ListenAndServe()
		if err != nil {
			logger.Error("could not start the server", "error", err)
			os.Exit(1)
		}
	}()

	// look for interrupts for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	logger.Info("shutting down the server", "received signal", sig)

	//gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	svr.Shutdown(ctx)
}
