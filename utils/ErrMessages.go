package utils

import "fmt"

var ErrUserAlreadyExists = fmt.Sprintf("User already exists with the given email")
var ErrUserNotFound = fmt.Sprintf("No user account exists with given email. Please sign up first")
var UserCreationFailed = fmt.Sprintf("Unable to create user.Please try again later")

var ErrArticleNotFound = fmt.Sprintf("Article not found")
var ErrCantUpdateOthersArticle = fmt.Sprintf("Only author can update the article!")
var ErrCantDeleteOthersArticle = fmt.Sprintf("Only author can delete the article!")
var ErrInvalidPageNumber = fmt.Sprintf("The requested page number is invalid.")
