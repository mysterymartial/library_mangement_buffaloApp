package Dto

import (
	"time"

	"github.com/gofrs/uuid"
)

type BookActionRequest struct {
	BookID uuid.UUID `json:"book_id"`
	Email  string    `json:"email"`
}

type BookActionResponse struct {
	ID         uuid.UUID  `json:"id"`
	UserID     uuid.UUID  `json:"userId"`
	BookID     uuid.UUID  `json:"bookId"`
	Status     string     `json:"status"`
	LoanDate   time.Time  `json:"loanDate"`
	ReturnDate *time.Time `json:"returnDate,omitempty"`
	Email      string     `json:"user_email"`
}
