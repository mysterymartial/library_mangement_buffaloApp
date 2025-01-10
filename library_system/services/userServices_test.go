package services

import (
	"errors"
	"library-system/Dto"
	"library-system/models"
	"library-system/repositories/mock"
	"log"

	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserServices_RegisterUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userRepo := &mock.MockUserRepo{}
		userService := UserServices{UserRepo: userRepo}

		request := Dto.UserRequest{
			Name:  "Aminat Usman",
			Email: "meenah20@gmail.com",
		}

		// Call the RegisterUser method
		user, err := userService.RegisterUser(request)

		// Add logging for debugging
		if err != nil {
			t.Errorf("Error during RegisterUser: %v", err)
		}
		log.Printf("User registered: %+v", user)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, user.Name, "aminat usman")
		assert.Equal(t, user.Email, "meenah20@gmail.com")
	})

	t.Run("Invalid Email", func(t *testing.T) {
		userRepo := &mock.MockUserRepo{}
		UserService := UserServices{UserRepo: userRepo}

		request := Dto.UserRequest{Name: "Aminat Usman", Email: "Invalid Email Address"}
		user, err := UserService.RegisterUser(request)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "Invalid Email Address", err.Error())
	})

	t.Run("Invalid Name", func(t *testing.T) {
		userRepo := &mock.MockUserRepo{}
		services := UserServices{UserRepo: userRepo}

		request := Dto.UserRequest{Name: "Aminat123", Email: "meenah20@gmail.com"}
		user, err := services.RegisterUser(request)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "Invalid Name", err.Error()) // This should fail due to name validation
	})

	t.Run("Error", func(t *testing.T) {
		userRepo := &mock.MockUserRepo{
			AddUserError: errors.New("failed to add user"),
		}
		service := UserServices{UserRepo: userRepo}

		request := Dto.UserRequest{Name: "Aminat Usman", Email: "meenah20@gmail.com"}
		user, err := service.RegisterUser(request)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "failed to add user", err.Error())
	})

	t.Run("duplicate email", func(t *testing.T) {
		userRepo := &mock.MockUserRepo{
			MockUser: []models.User{
				{
					Name:  "Existing User",
					Email: "meenah20@gmail.com",
				},
			},
		}
		userService := UserServices{UserRepo: userRepo}

		request := Dto.UserRequest{
			Name:  "Aminat Usman",
			Email: "meenah20@gmail.com",
		}

		user, err := userService.RegisterUser(request)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "email already registered", err.Error())
	})
}

func TestUserServices_CheckOutBook(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		bookID := uuid.Must(uuid.NewV4())
		userID := uuid.Must(uuid.NewV4())
		email := "meenah20@gmail.com"

		userRepo := &mock.MockUserRepo{
			MockUser: []models.User{
				{ID: userID, Name: "Aminat Usman", Email: email},
			},
		}
		bookRepo := &mock.MockBookRepository{
			MockBooks: []models.Book{
				{ID: bookID, Title: "Test Book", Status: "available"},
			},
		}
		loanRepo := &mock.MockLoanRepository{}
		service := UserServices{UserRepo: userRepo, LoanRepo: loanRepo, BookRepo: bookRepo}

		request := Dto.BookActionRequest{
			BookID: bookID,
			Email:  email,
		}

		response, err := service.CheckOutBook(request)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, "borrowed", response.Status)
		assert.Equal(t, email, response.Email)
	})

	t.Run("invalid email", func(t *testing.T) {
		bookID := uuid.Must(uuid.NewV4())
		service := UserServices{}

		request := Dto.BookActionRequest{
			BookID: bookID,
			Email:  "invalid-email",
		}

		response, err := service.CheckOutBook(request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "Invalid Email Address", err.Error())
	})

	t.Run("user not found", func(t *testing.T) {
		bookID := uuid.Must(uuid.NewV4())
		email := "nonexistent@example.com"

		userRepo := &mock.MockUserRepo{
			GetUserByEmailError: errors.New("user not found"),
		}
		bookRepo := &mock.MockBookRepository{
			MockBooks: []models.Book{
				{ID: bookID, Status: "available"},
			},
		}
		service := UserServices{
			UserRepo: userRepo,
			BookRepo: bookRepo,
		}

		request := Dto.BookActionRequest{
			BookID: bookID,
			Email:  email,
		}

		response, err := service.CheckOutBook(request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "User not found: user not found", err.Error())
	})

	t.Run("book not found", func(t *testing.T) {
		bookID := uuid.Must(uuid.NewV4())
		email := "meenah20@gmail.com"

		userRepo := &mock.MockUserRepo{
			MockUser: []models.User{
				{Email: email},
			},
		}
		bookRepo := &mock.MockBookRepository{
			GetBookByIDError: errors.New("book not found"),
		}
		service := UserServices{
			UserRepo: userRepo,
			BookRepo: bookRepo,
		}

		request := Dto.BookActionRequest{
			BookID: bookID,
			Email:  email,
		}

		response, err := service.CheckOutBook(request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "Book not found: book not found", err.Error())
	})

	t.Run("book already borrowed", func(t *testing.T) {
		bookID := uuid.Must(uuid.NewV4())
		email := "meenah20@gmail.com"

		userRepo := &mock.MockUserRepo{
			MockUser: []models.User{
				{Email: email},
			},
		}
		bookRepo := &mock.MockBookRepository{
			MockBooks: []models.Book{
				{ID: bookID, Status: "borrowed"}, // Book is already borrowed
			},
		}
		loanRepo := &mock.MockLoanRepository{
			MockLoans: []models.Loan{
				{BookID: bookID, Email: email, ReturnDate: nil},
			},
		}
		service := UserServices{
			UserRepo: userRepo,
			BookRepo: bookRepo,
			LoanRepo: loanRepo,
		}

		request := Dto.BookActionRequest{
			BookID: bookID,
			Email:  email,
		}

		response, err := service.CheckOutBook(request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "Book is not available for checkout", err.Error())
	})
}

func TestUserServices_ReturnBook(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		bookID := uuid.Must(uuid.NewV4())
		email := "meenah20@gmail.com"
		loanDate := time.Now().Add(-24 * time.Hour)

		loanRepo := &mock.MockLoanRepository{
			MockLoans: []models.Loan{
				{
					BookID:     bookID,
					Email:      email,
					LoanDate:   loanDate,
					ReturnDate: nil,
				},
			},
		}
		bookRepo := &mock.MockBookRepository{
			MockBooks: []models.Book{
				{ID: bookID, Status: "borrowed"},
			},
		}
		service := UserServices{
			LoanRepo: loanRepo,
			BookRepo: bookRepo,
		}

		request := Dto.BookActionRequest{
			BookID: bookID,
			Email:  email,
		}

		response, err := service.ReturnBook(request)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, "available", response.Status)
		assert.NotNil(t, response.ReturnDate)

		// Verify book status was updated
		updatedBook, _ := bookRepo.GetBookByID(bookID)
		assert.Equal(t, "available", updatedBook.Status)
	})
	t.Run("invalid email", func(t *testing.T) {
		bookID := uuid.Must(uuid.NewV4())
		request := Dto.BookActionRequest{
			BookID: bookID,
			Email:  "invalid-email",
		}
		service := UserServices{}

		response, err := service.ReturnBook(request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "Invalid Email Address", err.Error())
	})

	t.Run("no active loan found", func(t *testing.T) {
		bookID := uuid.Must(uuid.NewV4())
		email := "meenah20@gmail.com"

		loanRepo := &mock.MockLoanRepository{
			GetLoanByBookAndEmailError: errors.New("no loan found"),
		}
		service := UserServices{LoanRepo: loanRepo}

		request := Dto.BookActionRequest{
			BookID: bookID,
			Email:  email,
		}

		response, err := service.ReturnBook(request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "No active loan found for this book", err.Error())
	})

	t.Run("book already returned", func(t *testing.T) {
		bookID := uuid.Must(uuid.NewV4())
		email := "meenah20@gmail.com"
		loanDate := time.Now().Add(-24 * time.Hour)
		returnDate := time.Now()

		loanRepo := &mock.MockLoanRepository{
			MockLoans: []models.Loan{
				{
					ID:         uuid.Must(uuid.NewV4()),
					BookID:     bookID,
					Email:      email,
					LoanDate:   loanDate,
					ReturnDate: &returnDate,
				},
			},
		}
		bookRepo := &mock.MockBookRepository{
			MockBooks: []models.Book{
				{ID: bookID, Status: "available"},
			},
		}
		service := UserServices{
			LoanRepo: loanRepo,
			BookRepo: bookRepo,
		}

		request := Dto.BookActionRequest{
			BookID: bookID,
			Email:  email,
		}

		response, err := service.ReturnBook(request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "No active loan found for this book", err.Error())
	})
}

func TestUserServices_ReserveBook(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		bookID := uuid.Must(uuid.NewV4())
		userID := uuid.Must(uuid.NewV4())
		email := "meenah20@gmail.com"

		userRepo := &mock.MockUserRepo{
			MockUser: []models.User{
				{ID: userID, Email: email},
			},
		}
		bookRepo := &mock.MockBookRepository{
			MockBooks: []models.Book{
				{ID: bookID, Status: "available"},
			},
		}
		loanRepo := &mock.MockLoanRepository{
			MockLoans: []models.Loan{},
		}
		service := UserServices{
			UserRepo: userRepo,
			BookRepo: bookRepo,
			LoanRepo: loanRepo,
		}

		request := Dto.BookActionRequest{
			BookID: bookID,
			Email:  email,
		}

		response, err := service.ReserveBook(request)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, "reserved", response.Status)
		assert.Equal(t, email, response.Email)
	})

	t.Run("invalid email", func(t *testing.T) {
		bookID := uuid.Must(uuid.NewV4())
		service := UserServices{}

		request := Dto.BookActionRequest{
			BookID: bookID,
			Email:  "invalid-email",
		}

		response, err := service.ReserveBook(request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "Invalid Email Address", err.Error())
	})

	t.Run("user not found", func(t *testing.T) {
		bookID := uuid.Must(uuid.NewV4())
		email := "nonexistent@example.com"

		userRepo := &mock.MockUserRepo{
			GetUserByEmailError: errors.New("user not found"),
		}
		bookRepo := &mock.MockBookRepository{
			MockBooks: []models.Book{
				{ID: bookID, Status: "available"},
			},
		}
		service := UserServices{
			UserRepo: userRepo,
			BookRepo: bookRepo,
		}

		request := Dto.BookActionRequest{
			BookID: bookID,
			Email:  email,
		}

		response, err := service.ReserveBook(request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "User not found: user not found", err.Error())
	})

	t.Run("book not found", func(t *testing.T) {
		bookID := uuid.Must(uuid.NewV4())
		email := "meenah20@gmail.com"

		userRepo := &mock.MockUserRepo{
			MockUser: []models.User{
				{Email: email},
			},
		}
		bookRepo := &mock.MockBookRepository{
			GetBookByIDError: errors.New("book not found"),
		}
		service := UserServices{
			UserRepo: userRepo,
			BookRepo: bookRepo,
		}

		request := Dto.BookActionRequest{
			BookID: bookID,
			Email:  email,
		}

		response, err := service.ReserveBook(request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "Book not found")
	})

	t.Run("book not available", func(t *testing.T) {
		bookID := uuid.Must(uuid.NewV4())
		email := "meenah20@gmail.com"

		userRepo := &mock.MockUserRepo{
			MockUser: []models.User{
				{Email: email},
			},
		}
		bookRepo := &mock.MockBookRepository{
			MockBooks: []models.Book{
				{ID: bookID, Status: "borrowed"},
			},
		}
		service := UserServices{
			UserRepo: userRepo,
			BookRepo: bookRepo,
		}

		request := Dto.BookActionRequest{
			BookID: bookID,
			Email:  email,
		}

		response, err := service.ReserveBook(request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "Book is not available for reservation", err.Error())
	})
}
