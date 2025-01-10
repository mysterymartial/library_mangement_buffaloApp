package models

import (
	"errors"
	"github.com/gofrs/uuid"
	"strings"
	"time"
	"unicode"
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
		return errors.New("title is required")
	}
	if b.Author == "" {
		return errors.New("author is required")
	}
	if b.ISBN == "" {
		return errors.New("ISBN is required")
	}
	if b.Status == "" {
		return errors.New("status is required")
	}
	if !ValidateISBN(b.ISBN) {
		return errors.New("invalid ISBN format")
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
func ValidateISBN(isbn string) bool {
	isbn = strings.ReplaceAll(isbn, "-", "")
	isbn = strings.ReplaceAll(isbn, " ", "")

	if len(isbn) != 10 && len(isbn) != 13 {
		return false
	}

	for _, char := range isbn {
		if !unicode.IsDigit(char) && char != 'X' && char != 'x' {
			return false
		}
	}

	return true
}
