package services

import (
	"errors"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"library-system/Dto"
	"library-system/models"
	"library-system/repositories/mock"
	"testing"
	"time"
)

func setupTestService() (*BookServices, *mock.MockBookRepository) {
	mockRepo := mock.NewMockBookRepository()
	service := NewBookServices(mockRepo)
	return service, mockRepo
}

func TestBookServices_AddBook(t *testing.T) {
	service, _ := setupTestService()
	req := Dto.BookRequest{
		Title:  "Test Book",
		Author: "Test Author",
		ISBN:   "123456",
		Status: models.StatusAvailable,
	}

	book, err := service.AddBook(req)

	assert.NoError(t, err)
	assert.Equal(t, req.Title, book.Title)
	assert.Equal(t, req.Author, book.Author)
	assert.Equal(t, req.ISBN, book.ISBN)
	assert.Equal(t, models.StatusAvailable, book.Status)
}
func TestBookServices_AddBook_DuplicateISBN(t *testing.T) {

	existingBook := models.Book{
		ID:     uuid.Must(uuid.NewV4()),
		Title:  "Existing Book",
		Author: "Author",
		ISBN:   "123456",
	}
	mockRepo := &mock.MockBookRepository{
		MockBooks: []models.Book{existingBook},
	}
	service := NewBookServices(mockRepo)

	req := Dto.BookRequest{
		Title:  "New Book",
		Author: "New Author",
		ISBN:   "123456",
	}

	book, err := service.AddBook(req)

	assert.Error(t, err)
	assert.Nil(t, book)
	assert.Contains(t, err.Error(), "duplicate ISBN")
}

func TestBookServices_RemoveBook(t *testing.T) {
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

func TestBookServices_UpdateBook(t *testing.T) {
	bookID := uuid.Must(uuid.NewV4())
	mockBook := models.Book{
		ID:        bookID,
		Title:     "Original Title",
		Author:    "Original Author",
		ISBN:      "123456",
		Status:    models.StatusAvailable,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo := &mock.MockBookRepository{
		MockBooks: []models.Book{mockBook},
	}
	service := NewBookServices(mockRepo)

	request := Dto.BookRequest{
		ID:     bookID.String(),
		Title:  "Updated Title",
		Author: "Updated Author",
		ISBN:   "123456",
		Status: models.StatusBorrowed,
	}

	book, err := service.UpdateBook(request)

	assert.NoError(t, err)
	assert.NotNil(t, book)
	assert.Equal(t, bookID, book.ID)
	assert.Equal(t, request.Title, book.Title)
	assert.Equal(t, request.Author, book.Author)
	assert.Equal(t, request.Status, book.Status)
	assert.Equal(t, request.ISBN, book.ISBN)
}
func TestBookServices_UpdateBook_DuplicateISBN(t *testing.T) {

	book1ID := uuid.Must(uuid.NewV4())
	book2ID := uuid.Must(uuid.NewV4())
	existingBooks := []models.Book{
		{
			ID:     book1ID,
			Title:  "Book 1",
			Author: "Author 1",
			ISBN:   "123456",
		},
		{
			ID:     book2ID,
			Title:  "Book 2",
			Author: "Author 2",
			ISBN:   "789012",
		},
	}
	mockRepo := &mock.MockBookRepository{
		MockBooks: existingBooks,
	}
	service := NewBookServices(mockRepo)

	request := Dto.BookRequest{
		ID:     book2ID.String(),
		Title:  "Updated Book 2",
		Author: "Updated Author 2",
		ISBN:   "123456",
	}

	book, err := service.UpdateBook(request)

	assert.Error(t, err)
	assert.Nil(t, book)
	assert.Contains(t, err.Error(), "duplicate ISBN")
}

func TestBookServices_AddBookError(t *testing.T) {
	service, mockRepo := setupTestService()
	mockRepo.AddBookError = errors.New("failed to add book")

	request := Dto.BookRequest{
		Title:  "Test Book",
		Author: "Test Author",
		ISBN:   "123456",
		Status: models.StatusAvailable,
	}

	book, err := service.AddBook(request)

	assert.Error(t, err)
	assert.Nil(t, book)
	assert.Contains(t, err.Error(), "failed to add book")
}

func TestBookServices_RemoveBookError(t *testing.T) {
	service, mockRepo := setupTestService()
	bookID := uuid.Must(uuid.NewV4())
	mockRepo.GetBookByIDError = errors.New("book not found")

	book, err := service.RemoveBook(bookID)

	assert.Error(t, err)
	assert.Nil(t, book)
	assert.Contains(t, err.Error(), "failed to find book")
}

func TestBookServices_UpdateBookError(t *testing.T) {
	service, mockRepo := setupTestService()
	bookID := uuid.Must(uuid.NewV4())
	mockRepo.GetBookByIDError = errors.New("book not found")

	request := Dto.BookRequest{
		ID:     bookID.String(),
		Title:  "Updated Title",
		Author: "Updated Author",
		ISBN:   "123456",
		Status: models.StatusBorrowed,
	}

	book, err := service.UpdateBook(request)

	assert.Error(t, err)
	assert.Nil(t, book)
	assert.Contains(t, err.Error(), "failed to find book")
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
	// Create mock books
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
