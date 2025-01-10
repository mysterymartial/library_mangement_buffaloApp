package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"library-system/models"
	"log"
	"strings"
)

type BookRepository interface {
	AddBook(book *models.Book) error
	RemoveBook(bookID uuid.UUID) error
	GetBookByID(bookID uuid.UUID) (*models.Book, error)
	UpdateBook(book *models.Book) error
	SearchBook(query string) ([]*models.Book, error)
	GetBookByISBN(isbn string) (*models.Book, error)
	GetAllBooks() ([]*models.Book, error)
}

type BookRepositoryImpl struct {
	DB *pop.Connection
}

func NewBookRepository(db *pop.Connection) *BookRepositoryImpl {
	return &BookRepositoryImpl{DB: db}
}

func (r *BookRepositoryImpl) AddBook(book *models.Book) error {
	return r.DB.Transaction(func(tx *pop.Connection) error {
		existingBook := &models.Book{}
		err := tx.Where("isbn = ?", book.ISBN).First(existingBook)
		if err == nil {
			return fmt.Errorf("duplicate ISBN: book with ISBN %s already exists", book.ISBN)
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("error checking existing ISBN: %w", err)
		}

		if err := tx.Create(book); err != nil {
			return fmt.Errorf("error adding book: %w", err)
		}
		return nil
	})
}
func (r *BookRepositoryImpl) RemoveBook(bookID uuid.UUID) error {
	return r.DB.Transaction(func(tx *pop.Connection) error {
		book := &models.Book{}
		if err := tx.Find(book, bookID); err != nil {
			return fmt.Errorf("book with id %s not found", bookID)
		}
		return tx.Destroy(book)
	})
}

func (r *BookRepositoryImpl) GetBookByID(id uuid.UUID) (*models.Book, error) {
	book := &models.Book{}
	log.Printf("Attempting to fetch book with ID: %s", id.String())

	if err := r.DB.Find(book, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("book not found with id: %s", id)
		}
		return nil, fmt.Errorf("error finding book: %w", err)
	}

	log.Printf("Found book: %+v", book)
	return book, nil
}

func (r *BookRepositoryImpl) UpdateBook(book *models.Book) error {
	return r.DB.Transaction(func(tx *pop.Connection) error {
		if err := tx.Update(book); err != nil {
			return fmt.Errorf("error updating book: %w", err)
		}
		return nil
	})
}

func (r *BookRepositoryImpl) SearchBook(query string) ([]*models.Book, error) {
	var books []*models.Book
	query = strings.TrimSpace(query)
	q := r.DB.Where("title LIKE ? OR author LIKE ? OR isbn LIKE ?",
		"%"+query+"%", "%"+query+"%", "%"+query+"%")
	if err := q.All(&books); err != nil {
		return nil, fmt.Errorf("error searching books: %w", err)
	}
	return books, nil
}

func (r *BookRepositoryImpl) GetBookByISBN(isbn string) (*models.Book, error) {
	book := &models.Book{}
	err := r.DB.Where("isbn = ?", isbn).First(book)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("book with ISBN %s not found", isbn) // Match the test case expectation
		}
		return nil, fmt.Errorf("error finding book by ISBN: %w", err)
	}
	return book, nil
}

func (r *BookRepositoryImpl) GetAllBooks() ([]*models.Book, error) {
	var books []*models.Book
	if err := r.DB.All(&books); err != nil {
		return nil, fmt.Errorf("error fetching all books: %w", err)
	}
	return books, nil
}
