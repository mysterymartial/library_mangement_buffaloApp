package mock

import (
	"fmt"
	"sync"

	"github.com/gofrs/uuid"
	"library-system/models"
)

type MockBookRepository struct {
	sync.RWMutex
	MockBooks          []models.Book
	AddBookError       error
	RemoveBookError    error
	GetBookByIDError   error
	UpdateBookError    error
	SearchBookError    error
	GetBookByISBNError error
	GetAllBooksError   error
}

func NewMockBookRepository() *MockBookRepository {
	return &MockBookRepository{
		MockBooks: make([]models.Book, 0),
	}
}

func (r *MockBookRepository) AddBook(book *models.Book) error {
	r.Lock()
	defer r.Unlock()

	if r.AddBookError != nil {
		return r.AddBookError
	}

	r.MockBooks = append(r.MockBooks, *book)
	return nil
}

func (r *MockBookRepository) RemoveBook(bookID uuid.UUID) error {
	r.Lock()
	defer r.Unlock()

	if r.RemoveBookError != nil {
		return r.RemoveBookError
	}

	for i, book := range r.MockBooks {
		if bookID == book.ID {
			r.MockBooks = append(r.MockBooks[:i], r.MockBooks[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("book not found")
}

func (r *MockBookRepository) GetBookByID(bookID uuid.UUID) (*models.Book, error) {
	r.RLock()
	defer r.RUnlock()

	if r.GetBookByIDError != nil {
		return nil, r.GetBookByIDError
	}

	for _, book := range r.MockBooks {
		if book.ID == bookID {
			bookCopy := book
			return &bookCopy, nil
		}
	}
	return nil, fmt.Errorf("book not found")
}

func (r *MockBookRepository) UpdateBook(book *models.Book) error {
	r.Lock()
	defer r.Unlock()

	if r.UpdateBookError != nil {
		return r.UpdateBookError
	}

	for i, existingBook := range r.MockBooks {
		if existingBook.ID == book.ID {
			r.MockBooks[i] = *book
			return nil
		}
	}
	return fmt.Errorf("book not found")
}

func (r *MockBookRepository) SearchBook(query string) ([]*models.Book, error) {
	r.RLock()
	defer r.RUnlock()

	if r.SearchBookError != nil {
		return nil, r.SearchBookError
	}

	var results []*models.Book
	for _, book := range r.MockBooks {
		if book.Title == query || book.Author == query || book.ISBN == query {
			bookCopy := book
			results = append(results, &bookCopy)
		}
	}
	return results, nil
}

func (r *MockBookRepository) GetBookByISBN(isbn string) (*models.Book, error) {
	r.RLock()
	defer r.RUnlock()

	if r.GetBookByISBNError != nil {
		return nil, r.GetBookByISBNError
	}

	for _, book := range r.MockBooks {
		if book.ISBN == isbn {
			bookCopy := book
			return &bookCopy, nil
		}
	}
	return nil, nil
}

func (r *MockBookRepository) GetAllBooks() ([]*models.Book, error) {
	r.RLock()
	defer r.RUnlock()

	if r.GetAllBooksError != nil {
		return nil, r.GetAllBooksError
	}

	books := make([]*models.Book, len(r.MockBooks))
	for i := range r.MockBooks {
		bookCopy := r.MockBooks[i]
		books[i] = &bookCopy
	}
	return books, nil
}
