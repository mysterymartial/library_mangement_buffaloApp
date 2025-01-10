package services

import (
	"errors"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"library-system/Dto"
	"library-system/models"
	"library-system/repositories/mock"
)

func setupTestService() (*BookServices, *mock.MockBookRepository) {
	mockRepo := mock.NewMockBookRepository()
	service := NewBookServices(mockRepo)
	return service, mockRepo
}

func TestBookServices_TestThatYouCanSuccesfullyAddBook(t *testing.T) {
	service, _ := setupTestService()
	req := Dto.BookRequest{
		Title:  "Test Book",
		Author: "Test Author",
		ISBN:   "0-7475-3269-9", // Valid ISBN-10 format with hyphens
		Status: models.StatusAvailable,
	}

	book, err := service.AddBook(req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if book == nil {
		t.Fatal("Expected book to not be nil")
	}

	assert.Equal(t, req.Title, book.Title)
	assert.Equal(t, req.Author, book.Author)
	assert.Equal(t, req.ISBN, book.ISBN)
	assert.Equal(t, models.StatusAvailable, book.Status)
}

func TestBookServices_TestThatYouCanAddBookWithDuplicateISBN(t *testing.T) {
	existingBook := models.Book{
		ID:     uuid.Must(uuid.NewV4()),
		Title:  "Existing Book",
		Author: "Author",
		ISBN:   "978-3-16-148410-0",
	}
	mockRepo := &mock.MockBookRepository{
		MockBooks: []models.Book{existingBook},
	}
	service := NewBookServices(mockRepo)

	req := Dto.BookRequest{
		Title:  "New Book",
		Author: "New Author",
		ISBN:   "978-3-16-148410-0", // Duplicate ISBN
	}

	book, err := service.AddBook(req)

	assert.Error(t, err)
	assert.Nil(t, book)
	assert.Contains(t, err.Error(), "duplicate ISBN")
}

func TestBookServices_TestThatAddBookCanThrowAddBookError(t *testing.T) {
	service, mockRepo := setupTestService()
	mockRepo.AddBookError = errors.New("failed to add book")

	req := Dto.BookRequest{
		Title:  "Test Book",
		Author: "Test Author",
		ISBN:   "978-3-16-148410-0", // Valid ISBN-13 format
		Status: models.StatusAvailable,
	}

	book, err := service.AddBook(req)

	assert.Error(t, err)
	assert.Nil(t, book)
	assert.Contains(t, err.Error(), "failed to add book")
}

func TestBookServices_TestThatYouSuccessfullyRemoveBook(t *testing.T) {
	bookID := uuid.Must(uuid.NewV4())
	mockBook := models.Book{
		ID:        bookID,
		Title:     "Test Book",
		Author:    "Test Author",
		ISBN:      "123456",
		Status:    models.StatusAvailable,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo := &mock.MockBookRepository{
		MockBooks: []models.Book{mockBook},
	}
	service := NewBookServices(mockRepo)

	book, err := service.RemoveBook(bookID)

	assert.NoError(t, err)
	assert.NotNil(t, book)
	assert.Equal(t, bookID, book.ID)
	assert.Equal(t, mockBook.Title, book.Title)
	assert.Equal(t, mockBook.Author, book.Author)
	assert.Equal(t, mockBook.Status, book.Status)
	assert.Equal(t, mockBook.ISBN, book.ISBN)
}

func TestBookServices_TestThatCanThrowRemoveBookError(t *testing.T) {
	service, mockRepo := setupTestService()
	bookID := uuid.Must(uuid.NewV4())
	mockRepo.GetBookByIDError = errors.New("book not found")

	book, err := service.RemoveBook(bookID)

	assert.Error(t, err)
	assert.Nil(t, book)
	assert.Contains(t, err.Error(), "failed to find book")
}

func TestBookServices_TestYouCanUpdateBookByISBN(t *testing.T) {
	existingBook := models.Book{
		ID:        uuid.Must(uuid.NewV4()),
		Title:     "Original Title",
		Author:    "Original Author",
		ISBN:      "0-7475-3269-9",
		Status:    models.StatusAvailable,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo := &mock.MockBookRepository{
		MockBooks: []models.Book{existingBook},
	}
	service := NewBookServices(mockRepo)

	req := Dto.BookRequest{
		Title:  "Updated Title",
		Author: "Updated Author",
		ISBN:   "0-7475-3269-9",
		Status: models.StatusBorrowed,
	}

	book, err := service.UpdateBookByISBN(req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if book == nil {
		t.Fatal("Expected book to not be nil")
	}

	assert.Equal(t, existingBook.ID, book.ID)
	assert.Equal(t, req.Title, book.Title)
	assert.Equal(t, req.Author, book.Author)
	assert.Equal(t, req.Status, book.Status)
	assert.Equal(t, req.ISBN, book.ISBN)
}

func TestBookServices_TestUpdateBookByISBN_NotFound(t *testing.T) {
	mockRepo := &mock.MockBookRepository{
		MockBooks: []models.Book{},
	}
	service := NewBookServices(mockRepo)

	req := Dto.BookRequest{
		Title:  "Non-existent Book",
		Author: "Non-existent Author",
		ISBN:   "0-7475-3269-9", // ISBN that does not exist in the repository
		Status: models.StatusAvailable,
	}

	book, err := service.UpdateBookByISBN(req)

	assert.Error(t, err)
	assert.Nil(t, book)
	assert.Contains(t, err.Error(), "book with ISBN 0-7475-3269-9 not found")
}

func TestBookServices_UpdateBookByISBN_DatabaseError(t *testing.T) {
	existingBook := models.Book{
		ID:     uuid.Must(uuid.NewV4()),
		Title:  "Original Title",
		Author: "Original Author",
		ISBN:   "123456",
		Status: models.StatusAvailable,
	}

	mockRepo := &mock.MockBookRepository{
		MockBooks:       []models.Book{existingBook},
		UpdateBookError: errors.New("database error"),
	}
	service := NewBookServices(mockRepo)

	req := Dto.BookRequest{
		Title:  "Updated Title",
		Author: "Updated Author",
		ISBN:   "123456", // Valid ISBN
		Status: models.StatusAvailable,
	}

	book, err := service.UpdateBookByISBN(req)

	assert.Error(t, err)
	assert.Nil(t, book)
	assert.Contains(t, err.Error(), "validation")
}

func TestBookServices_SearchBook(t *testing.T) {
	service, mockRepo := setupTestService()
	mockRepo.SearchBookError = errors.New("search failed")

	books, err := service.SearchBook("Test Book")

	assert.Error(t, err)
	assert.Nil(t, books)
	assert.Contains(t, err.Error(), "failed to search books")
}

func TestBookServices_GetBookByID(t *testing.T) {
	bookID := uuid.Must(uuid.NewV4())
	mockBook := models.Book{
		ID:        bookID,
		Title:     "Test Book",
		Author:    "Test Author",
		ISBN:      "123456",
		Status:    models.StatusAvailable,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo := &mock.MockBookRepository{
		MockBooks: []models.Book{mockBook},
	}
	service := NewBookServices(mockRepo)

	book, err := service.GetBookByID(bookID)

	assert.NoError(t, err)
	assert.NotNil(t, book)
	assert.Equal(t, mockBook.ID, book.ID)
	assert.Equal(t, mockBook.Title, book.Title)
}

func TestBookServices_GetBookByISBN_Error(t *testing.T) {
	_, mockRepo := setupTestService()
	mockRepo.GetBookByISBNError = errors.New("database error")

	existingBook, err := mockRepo.GetBookByISBN("123456")
	assert.Error(t, err)
	assert.Nil(t, existingBook)
	assert.Contains(t, err.Error(), "database error")
}

func TestBookServices_GetAllBooks(t *testing.T) {
	mockBooks := []models.Book{
		{
			ID:     uuid.Must(uuid.NewV4()),
			Title:  "Book 1",
			Author: "Author 1",
			ISBN:   "123456",
		},
		{
			ID:     uuid.Must(uuid.NewV4()),
			Title:  "Book 2",
			Author: "Author 2",
			ISBN:   "789012",
		},
	}

	mockRepo := &mock.MockBookRepository{
		MockBooks: mockBooks,
	}
	service := NewBookServices(mockRepo)

	books, err := service.GetAllBooks()

	assert.NoError(t, err)
	assert.NotNil(t, books)
	assert.Equal(t, len(mockBooks), len(books))
	assert.Equal(t, mockBooks[0].Title, books[0].Title)
	assert.Equal(t, mockBooks[1].Title, books[1].Title)
}
