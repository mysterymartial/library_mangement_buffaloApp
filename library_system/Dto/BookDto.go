package Dto

import "github.com/gobuffalo/uuid"

type BookRequest struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	ISBN      string `json:"isbn"`
	Status    string `json:"status"`
	UserToken string `json:"-"`
}

type BookResponse struct {
	ID     uuid.UUID `json:"id"`
	Title  string    `json:"title"`
	Author string    `json:"author"`
	ISBN   string    `json:"isbn"`
	Status string    `json:"status"`
}
