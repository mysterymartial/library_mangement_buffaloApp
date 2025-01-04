package Dto

import (
	"github.com/gofrs/uuid"
	"time"
)

type LoanRequest struct {
	BookID uuid.UUID `json:"book_id" validate:"required"`
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

type LoanResponse struct {
	ID         uuid.UUID  `json:"id"`
	BookID     uuid.UUID  `json:"book_id"`
	UserID     uuid.UUID  `json:"user_id"`
	LoanDate   time.Time  `json:"loan_date"`
	ReturnDate *time.Time `json:"return_date,omitempty"`
	UserName   string     `json:"user_name,omitempty"`
}
