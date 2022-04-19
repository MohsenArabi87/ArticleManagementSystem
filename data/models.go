package data

import "time"

// User is the data type for user object
type User struct {
	Email     string `json:"email" validate:"required" `
	Password  string `json:"password" validate:"required"`
	TokenHash string `json:"tokenhash"`
}

//Article is the data type for article object
type Article struct {
	ID        string    `json:"ID"`
	Title     string    `json:"title" validate:"required" `
	Content   string    `json:"content" `
	Tags      []string  `json:"tags" `
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
