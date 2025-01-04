package controllers

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gofrs/uuid"
	"library-system/Dto"
	"library-system/services"
	"net/http"
	"strings"
)

var r *render.Engine

// Initialize render engine
func init() {
	r = render.New(render.Options{})
}

type BookController struct {
	BookService *services.BookServices
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

func NewBookController(bookService *services.BookServices) *BookController {
	return &BookController{BookService: bookService}
}

// AddBook handles adding a new book to the system
func (bc *BookController) AddBook(c buffalo.Context) error {
	var request Dto.BookRequest
	// Bind the request data into the BookRequest DTO
	if err := c.Bind(&request); err != nil {
		return c.Render(http.StatusBadRequest, r.JSON(ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		}))
	}

	// Validate ISBN is provided
	if strings.TrimSpace(request.ISBN) == "" {
		return c.Render(http.StatusBadRequest, r.JSON(ErrorResponse{
			Error: "ISBN is required",
		}))
	}

	// Use the BookService to add the book
	book, err := bc.BookService.AddBook(request)
	if err != nil {
		// Handle duplicate ISBN error or other issues
		if strings.Contains(err.Error(), "duplicate ISBN") {
			return c.Render(http.StatusConflict, r.JSON(ErrorResponse{
				Error: err.Error(),
			}))
		}
		return handleError(c, err)
	}

	// Return the created book with its ID
	return c.Render(http.StatusCreated, r.JSON(book))
}

// UpdateBook handles updating an existing book in the system
func (bc *BookController) UpdateBook(c buffalo.Context) error {
	var request Dto.BookRequest
	// Bind the request data into the BookRequest DTO
	if err := c.Bind(&request); err != nil {
		return c.Render(http.StatusBadRequest, r.JSON(ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		}))
	}

	// Validate ISBN is provided
	if strings.TrimSpace(request.ISBN) == "" {
		return c.Render(http.StatusBadRequest, r.JSON(ErrorResponse{
			Error: "ISBN is required",
		}))
	}

	// Use the BookService to update the book
	book, err := bc.BookService.UpdateBook(request)
	if err != nil {
		// Handle duplicate ISBN error or other issues
		if strings.Contains(err.Error(), "duplicate ISBN") {
			return c.Render(http.StatusConflict, r.JSON(ErrorResponse{
				Error: err.Error(),
			}))
		}
		return handleError(c, err)
	}

	// Return the updated book
	return c.Render(http.StatusOK, r.JSON(book))
}

// RemoveBook handles removing a book from the system by its ID
func (bc *BookController) RemoveBook(c buffalo.Context) error {
	bookID, err := parseUUID(c.Param("id")) // Get ID from the URL path
	if err != nil {
		return c.Render(http.StatusBadRequest, r.JSON(ErrorResponse{
			Error:   "Invalid book ID format",
			Details: err.Error(),
		}))
	}

	// Use the BookService to remove the book
	book, err := bc.BookService.RemoveBook(bookID)
	if err != nil {
		return handleError(c, err)
	}

	// Return the removed book details (or could just return a success message)
	return c.Render(http.StatusOK, r.JSON(book))
}

// GetBookByID retrieves a single book by its ID
func (bc *BookController) GetBookByID(c buffalo.Context) error {
	bookID, err := parseUUID(c.Param("id"))
	if err != nil {
		return c.Render(http.StatusBadRequest, r.JSON(ErrorResponse{
			Error:   "Invalid book ID format",
			Details: err.Error(),
		}))
	}

	// Use the BookService to fetch the book by ID
	book, err := bc.BookService.GetBookByID(bookID)
	if err != nil {
		return handleError(c, err)
	}

	// Return the book details
	return c.Render(http.StatusOK, r.JSON(book))
}

// SearchBook searches for books based on a query string
func (bc *BookController) SearchBook(c buffalo.Context) error {
	query := c.Param("query")
	if strings.TrimSpace(query) == "" {
		return c.Render(http.StatusBadRequest, r.JSON(ErrorResponse{
			Error: "Search query cannot be empty",
		}))
	}

	// Use the BookService to search for books
	books, err := bc.BookService.SearchBook(query)
	if err != nil {
		return handleError(c, err)
	}

	// Return the search results
	return c.Render(http.StatusOK, r.JSON(books))
}

// GetAllBooks retrieves all books in the system
func (bc *BookController) GetAllBooks(c buffalo.Context) error {
	books, err := bc.BookService.GetAllBooks()
	if err != nil {
		return c.Render(http.StatusInternalServerError, r.JSON(ErrorResponse{
			Error:   "Failed to fetch books",
			Details: err.Error(),
		}))
	}

	// Return the list of all books
	return c.Render(http.StatusOK, r.JSON(books))
}

// Helper function to parse UUID
func parseUUID(id string) (uuid.UUID, error) {
	return uuid.FromString(id)
}

// Global error handling
func handleError(c buffalo.Context, err error) error {
	errLower := strings.ToLower(err.Error())

	// Error handling based on the error type
	switch {
	case strings.Contains(errLower, "not found"):
		return c.Render(http.StatusNotFound, r.JSON(ErrorResponse{
			Error: err.Error(),
		}))
	case strings.Contains(errLower, "validation"):
		return c.Render(http.StatusBadRequest, r.JSON(ErrorResponse{
			Error: err.Error(),
		}))
	case strings.Contains(errLower, "duplicate isbn"):
		return c.Render(http.StatusConflict, r.JSON(ErrorResponse{
			Error: err.Error(),
		}))
	default:
		return c.Render(http.StatusInternalServerError, r.JSON(ErrorResponse{
			Error:   "Internal server error",
			Details: err.Error(),
		}))
	}
}
