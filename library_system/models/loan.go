package models

import (
	"errors"
	"github.com/gofrs/uuid"
	"time"
)

type Loan struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	BookID     uuid.UUID  `json:"book_id" db:"book_id"`
	UserID     uuid.UUID  `json:"user_id" db:"user_id"`
	UserName   string     `json:"user_name" db:"user_name"`
	LoanDate   time.Time  `json:"loan_date" db:"loan_date"`
	ReturnDate *time.Time `json:"return_date" db:"return_date"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
}

func (l *Loan) Validate() error {
	if l.BookID == uuid.Nil {
		return errors.New("book ID is required")
	}
	if l.UserID == uuid.Nil {
		return errors.New("user ID is required")
	}
	return nil
}
