package models

import (
	"errors"
	"github.com/gofrs/uuid"
	"time"
)

const (
	StatusAvailable = "available"
	StatusBorrowed  = "borrowed"
	StatusReserved  = "reserved"
)

type Book struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Author    string    `json:"author" db:"author"`
	ISBN      string    `json:"isbn" db:"isbn"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (b *Book) Validate() error {
	if b.Title == "" {
		return errors.New("title cannot be empty")
	}
	if b.Author == "" {
		return errors.New("author cannot be empty")
	}
	if b.ISBN == "" {
		return errors.New("ISBN cannot be empty")
	}
	if b.Status == "" {
		return errors.New("status cannot be empty")
	}

	validStatuses := map[string]bool{
		StatusAvailable: true,
		StatusBorrowed:  true,
		StatusReserved:  true,
	}

	if !validStatuses[b.Status] {
		return errors.New("invalid status value")
	}

	return nil
}
