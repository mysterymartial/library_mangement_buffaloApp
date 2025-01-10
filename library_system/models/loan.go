package models

import (
	"errors"
	"github.com/gofrs/uuid"
	"time"
)

type Loan struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	BookID     uuid.UUID  `json:"book_id" db:"book_id"`
	Email      string     `json:"email" db:"email"`
	UserID     uuid.UUID  `json:"user_id" db:"user_id"`
	LoanDate   time.Time  `json:"loan_date" db:"loan_date"`
	ReturnDate *time.Time `json:"return_date" db:"return_date"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
}

func (l *Loan) Validate() error {
	if l.BookID == uuid.Nil {
		return errors.New("book ID is required")
	}
	if l.Email == "" {
		return errors.New("email is required")
	}
	return nil
}

func (l *Loan) BeforeCreate() error {
	if l.ID == uuid.Nil {
		id, err := uuid.NewV4()
		if err != nil {
			return err
		}
		l.ID = id
	}

	if l.UserID == uuid.Nil {
		l.UserID = uuid.Nil
	}

	now := time.Now()
	if l.CreatedAt.IsZero() {
		l.CreatedAt = now
	}
	if l.UpdatedAt.IsZero() {
		l.UpdatedAt = now
	}
	if l.LoanDate.IsZero() {
		l.LoanDate = now
	}

	return nil
}
