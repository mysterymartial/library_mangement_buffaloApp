package Dto

import "github.com/gobuffalo/uuid"

// BookRequest is the DTO for creating/updating a book
type BookRequest struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	ISBN      string `json:"isbn"`
	Status    string `json:"status"`
	UserToken string `json:"-"`
}

// BookResponse is the DTO for returning a book
type BookResponse struct {
	ID     uuid.UUID `json:"id"`
	Title  string    `json:"title"`
	Author string    `json:"author"`
	ISBN   string    `json:"isbn"`
	Status string    `json:"status"`
}
