package controllers

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gofrs/uuid"
	"library-system/Dto"
	"library-system/models"
	"library-system/services"
	"net/http"
	"strings"
)

var r *render.Engine

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

func (bc *BookController) AddBook(c buffalo.Context) error {
	var request Dto.BookRequest
	if err := c.Bind(&request); err != nil {
		return c.Render(http.StatusBadRequest, r.JSON(ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		}))
	}

	if strings.TrimSpace(request.ISBN) == "" {
		return c.Render(http.StatusBadRequest, r.JSON(ErrorResponse{
			Error: "ISBN is required",
		}))
	}

	book, err := bc.BookService.AddBook(request)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate ISBN") {
			return c.Render(http.StatusConflict, r.JSON(ErrorResponse{
				Error: err.Error(),
			}))
		}
		return handleError(c, err)
	}

	return c.Render(http.StatusCreated, r.JSON(book))
}

func (bc *BookController) RemoveBook(c buffalo.Context) error {
	bookID, err := parseUUID(c.Param("id")) // Get ID from the URL path
	if err != nil {
		return c.Render(http.StatusBadRequest, r.JSON(ErrorResponse{
			Error:   "Invalid book ID format",
			Details: err.Error(),
		}))
	}

	book, err := bc.BookService.RemoveBook(bookID)
	if err != nil {
		return handleError(c, err)
	}

	return c.Render(http.StatusOK, r.JSON(book))
}

func (bc *BookController) UpdateBook(c buffalo.Context) error {
	var request Dto.BookRequest
	if err := c.Bind(&request); err != nil {
		return c.Render(http.StatusBadRequest, r.JSON(ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		}))
	}

	if strings.TrimSpace(request.ISBN) == "" {
		return c.Render(http.StatusBadRequest, r.JSON(ErrorResponse{
			Error: "ISBN is required",
		}))
	}

	request.Status = models.StatusBorrowed

	book, err := bc.BookService.UpdateBookByISBN(request)
	if err != nil {
		return handleError(c, err)
	}

	return c.Render(http.StatusOK, r.JSON(book))
}

func (bc *BookController) GetBookByID(c buffalo.Context) error {
	bookID, err := parseUUID(c.Param("id"))
	if err != nil {
		return c.Render(http.StatusBadRequest, r.JSON(ErrorResponse{
			Error:   "Invalid book ID format",
			Details: err.Error(),
		}))
	}

	book, err := bc.BookService.GetBookByID(bookID)
	if err != nil {
		return handleError(c, err)
	}

	return c.Render(http.StatusOK, r.JSON(book))
}

func (bc *BookController) SearchBook(c buffalo.Context) error {
	query := c.Param("query")
	if strings.TrimSpace(query) == "" {
		return c.Render(http.StatusBadRequest, r.JSON(ErrorResponse{
			Error: "Search query cannot be empty",
		}))
	}

	books, err := bc.BookService.SearchBook(query)
	if err != nil {
		return handleError(c, err)
	}

	return c.Render(http.StatusOK, r.JSON(books))
}

func (bc *BookController) GetAllBooks(c buffalo.Context) error {
	books, err := bc.BookService.GetAllBooks()
	if err != nil {
		return c.Render(http.StatusInternalServerError, r.JSON(ErrorResponse{
			Error:   "Failed to fetch books",
			Details: err.Error(),
		}))
	}

	return c.Render(http.StatusOK, r.JSON(books))
}

func parseUUID(id string) (uuid.UUID, error) {
	return uuid.FromString(id)
}

func handleError(c buffalo.Context, err error) error {
	errLower := strings.ToLower(err.Error())

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
