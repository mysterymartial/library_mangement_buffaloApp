package services

import (
	"fmt"
	"github.com/gofrs/uuid"
	"library-system/Dto"
	"library-system/models"
	"library-system/repositories/repository"
	"regexp"
	"strings"
	"time"
)

type BookServices struct {
	BookRepo repository.BookRepository
}

func NewBookServices(bookRepo repository.BookRepository) *BookServices {
	return &BookServices{
		BookRepo: bookRepo,
	}
}

func (s *BookServices) AddBook(req Dto.BookRequest) (*Dto.BookResponse, error) {
	if strings.TrimSpace(req.ISBN) != "" {
		if !isValidISBN(req.ISBN) {
			return nil, fmt.Errorf("invalid ISBN format")
		}

		existingBook, err := s.BookRepo.GetBookByISBN(req.ISBN)
		if err != nil {
			if !strings.Contains(err.Error(), "not found") {
				return nil, fmt.Errorf("error checking ISBN: %w", err)
			}
		}
		if existingBook != nil {
			return nil, fmt.Errorf("duplicate ISBN: book with ISBN %s already exists", req.ISBN)
		}
	}

	book := &models.Book{
		ID:        uuid.Must(uuid.NewV4()),
		Title:     req.Title,
		Author:    req.Author,
		ISBN:      req.ISBN,
		Status:    models.StatusAvailable,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := book.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	if err := s.BookRepo.AddBook(book); err != nil {
		return nil, fmt.Errorf("failed to add book: %w", err)
	}

	return mapBookToResponse(book), nil
}

func (s *BookServices) RemoveBook(bookID uuid.UUID) (*Dto.BookResponse, error) {
	book, err := s.BookRepo.GetBookByID(bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to find book: %w", err)
	}

	if err := s.BookRepo.RemoveBook(bookID); err != nil {
		return nil, fmt.Errorf("failed to remove book: %w", err)
	}

	return mapBookToResponse(book), nil
}

func (s *BookServices) UpdateBookByISBN(request Dto.BookRequest) (*Dto.BookResponse, error) {

	existingBook, err := s.BookRepo.GetBookByISBN(request.ISBN)
	if err != nil {
		return nil, fmt.Errorf("book with ISBN %s not found", request.ISBN)
	}

	if existingBook == nil {
		return nil, fmt.Errorf("book with ISBN %s not found", request.ISBN)
	}

	if existingBook.ISBN != request.ISBN {
		return nil, fmt.Errorf("cannot update ISBN")
	}
	existingBook.Title = request.Title
	existingBook.Author = request.Author
	existingBook.Status = request.Status
	existingBook.UpdatedAt = time.Now()

	if err := existingBook.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	if err := s.BookRepo.UpdateBook(existingBook); err != nil {
		return nil, fmt.Errorf("failed to update book: %w", err)
	}

	return mapBookToResponse(existingBook), nil
}

func (s *BookServices) SearchBook(query string) ([]Dto.BookResponse, error) {
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	books, err := s.BookRepo.SearchBook(query)
	if err != nil {
		return nil, fmt.Errorf("failed to search books: %w", err)
	}

	return mapBooksToResponses(books), nil
}

func mapBookToResponse(book *models.Book) *Dto.BookResponse {
	if book == nil {
		return nil
	}
	return &Dto.BookResponse{
		ID:     book.ID,
		Title:  book.Title,
		Author: book.Author,
		ISBN:   book.ISBN,
		Status: book.Status,
	}
}

func mapBooksToResponses(books []*models.Book) []Dto.BookResponse {
	responses := make([]Dto.BookResponse, 0, len(books))
	for _, book := range books {
		if response := mapBookToResponse(book); response != nil {
			responses = append(responses, *response)
		}
	}
	return responses
}

func (s *BookServices) GetAllBooks() ([]Dto.BookResponse, error) {
	books, err := s.BookRepo.GetAllBooks()
	if err != nil {
		return nil, fmt.Errorf("failed to get all books: %w", err)
	}

	return mapBooksToResponses(books), nil
}

func (s *BookServices) GetBookByID(bookID uuid.UUID) (*Dto.BookResponse, error) {
	book, err := s.BookRepo.GetBookByID(bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to find book: %w", err)
	}

	return mapBookToResponse(book), nil
}

func isValidISBN(isbn string) bool {
	isbn = strings.ReplaceAll(isbn, "-", "")
	isbn10Regex := regexp.MustCompile(`^\d{9}[\dXx]$`)
	isbn13Regex := regexp.MustCompile(`^\d{13}$`)
	return isbn10Regex.MatchString(isbn) || isbn13Regex.MatchString(isbn)
}
