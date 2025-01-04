package Dto

import (
	"time"

	"github.com/gofrs/uuid"
)

type BookActionRequest struct {
	BookID   uuid.UUID `json:"book_id"`
	UserID   uuid.UUID `json:"user_id,omitempty"`
	UserName string    `json:"user_name"`
}

type BookActionResponse struct {
	ID         uuid.UUID  `json:"id"`
	UserID     uuid.UUID  `json:"userId"`
	BookID     uuid.UUID  `json:"bookId"`
	Status     string     `json:"status"`
	LoanDate   time.Time  `json:"loanDate"`
	ReturnDate *time.Time `json:"returnDate,omitempty"`
	UserName   string     `json:"userName"`
}
