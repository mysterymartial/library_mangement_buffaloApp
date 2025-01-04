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
		bookID, _ := uuid.NewV4()
		userID, _ := uuid.NewV4()
		userRepo := &mock.MockUserRepo{
			MockUser: []models.User{
				{ID: userID, Name: "Aminat Usman", Email: "meenah20@gmail.com"},
			},
		}
		bookRepo := &mock.MockBookRepository{
			MockBooks: []models.Book{
				{ID: bookID, Title: "Test Book", Author: "Test Author", ISBN: "123456", Status: "available"},
			},
		}
		loanRepo := &mock.MockLoanRepository{
			MockLoans: []models.Loan{},
		}
		service := UserServices{UserRepo: userRepo, LoanRepo: loanRepo, BookRepo: bookRepo}

		request := Dto.BookActionRequest{
			UserID:   userID,
			BookID:   bookID,
			UserName: "Aminat Usman",
		}

		action, err := service.CheckOutBook(request)

		assert.NoError(t, err)
		assert.NotNil(t, action)
		assert.Equal(t, "borrowed", action.Status)
	})

	t.Run("User Not Found", func(t *testing.T) {
		bookID, _ := uuid.NewV4()
		userID, _ := uuid.NewV4()
		userRepo := &mock.MockUserRepo{
			GetUserByIDError: errors.New("User Not Found"),
		}
		bookRepo := &mock.MockBookRepository{}
		loanRepo := &mock.MockLoanRepository{}
		service := UserServices{userRepo, loanRepo, bookRepo}

		request := Dto.BookActionRequest{UserID: userID, BookID: bookID, UserName: "Aminat Usman"}
		action, err := service.CheckOutBook(request)
		assert.Error(t, err)
		assert.Nil(t, action)
		assert.Equal(t, "User not found", err.Error())

	})
	t.Run("UserName Mismatch", func(t *testing.T) {
		bookID, _ := uuid.NewV4()
		userID, _ := uuid.NewV4()
		userRepo := &mock.MockUserRepo{
			MockUser: []models.User{
				{ID: userID, Name: "Aminat Usman", Email: "meenah20@gmail.com"},
			},
		}
		bookRepo := &mock.MockBookRepository{
			MockBooks: []models.Book{
				{ID: bookID, Title: "Test Book", Status: "Available"},
			},
		}
		loanRepo := &mock.MockLoanRepository{}
		service := UserServices{userRepo, loanRepo, bookRepo}

		request := Dto.BookActionRequest{UserID: userID, BookID: bookID, UserName: "Aminat Agbaosi"}
		action, err := service.CheckOutBook(request)
		assert.Error(t, err)
		assert.Nil(t, action)
		assert.Equal(t, "Invalid User Name", err.Error())
	})

}
func TestUserServices_ReturnBook(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		bookID, _ := uuid.NewV4()
		userID, _ := uuid.NewV4()
		loanID, _ := uuid.NewV4()
		loanDate := time.Now().AddDate(0, 0, -7)
		loanRepo := &mock.MockLoanRepository{
			MockLoans: []models.Loan{
				{ID: loanID, BookID: bookID, UserID: userID, LoanDate: loanDate, UserName: "Aminat Usman"},
			},
		}
		userRepo := &mock.MockUserRepo{
			MockUser: []models.User{
				{ID: userID, Name: "Aminat Usman", Email: "meenah20@gmail.com"},
			},
		}
		bookRepo := &mock.MockBookRepository{
			MockBooks: []models.Book{
				{ID: bookID, Title: "Test Book", Status: "Available"},
			},
		}
		service := UserServices{userRepo, loanRepo, bookRepo}
		request := Dto.BookActionRequest{UserID: userID, BookID: bookID, UserName: "Aminat Usman"}
		action, err := service.ReturnBook(request)

		// Check for errors
		assert.NoError(t, err)
		assert.NotNil(t, action)

		// Check the returned status and return date
		assert.Equal(t, "Returned", action.Status)

		// Ensure the ReturnDate is within a second of the current time
		if action.ReturnDate != nil {
			assert.WithinDuration(t, time.Now(), *action.ReturnDate, time.Second)
		} else {
			t.Errorf("Expected non-nil ReturnDate, got nil")
		}

		// Check the user name
		assert.Equal(t, "Aminat Usman", action.UserName)
	})

	t.Run("User Not Found", func(t *testing.T) {
		bookID, _ := uuid.NewV4()
		userID, _ := uuid.NewV4()
		userRepo := &mock.MockUserRepo{
			GetUserByIDError: errors.New("User Not Found"),
		}
		bookRepo := &mock.MockBookRepository{}
		loanRepo := &mock.MockLoanRepository{}
		service := UserServices{userRepo, loanRepo, bookRepo}

		request := Dto.BookActionRequest{UserID: userID, BookID: bookID, UserName: "Aminat Usman"}
		action, err := service.ReturnBook(request)
		assert.Error(t, err)
		assert.Nil(t, action)
		assert.Equal(t, "User Not Found or Name Doesn't Match", err.Error())
	})

	t.Run("UserName Mismatch", func(t *testing.T) {
		bookID, _ := uuid.NewV4()
		userID, _ := uuid.NewV4()
		loanID, _ := uuid.NewV4()
		loanDate := time.Now().AddDate(0, 0, -7)
		loanRepo := &mock.MockLoanRepository{
			MockLoans: []models.Loan{
				{ID: loanID, BookID: bookID, UserID: userID, LoanDate: loanDate, UserName: "Aminat Usman"},
			},
		}
		userRepo := &mock.MockUserRepo{
			MockUser: []models.User{
				{ID: userID, Name: "Aminat Usman", Email: "meenah20@gmail.com"},
			},
		}
		bookRepo := &mock.MockBookRepository{
			MockBooks: []models.Book{
				{ID: bookID, Title: "Test Book", Author: "Test Author", Status: "Available"},
			},
		}
		service := UserServices{userRepo, loanRepo, bookRepo}

		request := Dto.BookActionRequest{UserID: userID, BookID: bookID, UserName: "Aminat Agbaosi"}
		action, err := service.ReturnBook(request)
		assert.Error(t, err)
		assert.Nil(t, action)
		assert.Equal(t, "User Not Found or Name Doesn't Match", err.Error())
	})

	t.Run("Loan Not Found", func(t *testing.T) {
		bookID, _ := uuid.NewV4()
		userID, _ := uuid.NewV4()
		userRepo := &mock.MockUserRepo{
			MockUser: []models.User{
				{ID: userID, Name: "Aminat Usman", Email: "meenah20@gmail.com"},
			},
		}
		bookRepo := &mock.MockBookRepository{}
		loanRepo := &mock.MockLoanRepository{
			GetLoanByBookAndUserError: errors.New("Loan Not Found"),
		}
		service := UserServices{userRepo, loanRepo, bookRepo}
		request := Dto.BookActionRequest{UserID: userID, BookID: bookID, UserName: "Aminat Usman"}
		action, err := service.ReturnBook(request)
		assert.Error(t, err)
		assert.Nil(t, action)
		assert.Equal(t, "Loan Record Not Found", err.Error())

	})
}

func TestUserServices_ReserveBook(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		bookID, _ := uuid.NewV4()
		userID, _ := uuid.NewV4()
		userRepo := &mock.MockUserRepo{
			MockUser: []models.User{
				{ID: userID, Name: "Aminat Usman", Email: "meenah20@gmail.com"},
			},
		}
		bookRepo := &mock.MockBookRepository{
			MockBooks: []models.Book{
				{ID: bookID, Title: "Test Book", Author: "Test Author", Status: "Available"},
			},
		}
		service := UserServices{UserRepo: userRepo, BookRepo: bookRepo}

		request := Dto.BookActionRequest{UserID: userID, BookID: bookID, UserName: "Aminat Usman"}
		action, err := service.ReserveBook(request)
		assert.NoError(t, err)
		assert.NotNil(t, action)
		assert.Equal(t, "Reserved", action.Status)
		assert.Equal(t, "Aminat Usman", action.UserName)
	})
	t.Run("User Not Found", func(t *testing.T) {
		bookID, _ := uuid.NewV4()
		userID, _ := uuid.NewV4()
		userRepo := &mock.MockUserRepo{
			GetUserByIDError: errors.New("User Not Found"),
		}
		bookRepo := &mock.MockBookRepository{}
		service := UserServices{UserRepo: userRepo, BookRepo: bookRepo}
		request := Dto.BookActionRequest{UserID: userID, BookID: bookID, UserName: "Aminat Usman"}
		action, err := service.ReturnBook(request)
		assert.Error(t, err)
		assert.Nil(t, action)
		assert.Equal(t, "User Not Found or Name Doesn't Match", err.Error())
	})
	t.Run("User Mismatch", func(t *testing.T) {
		bookID, _ := uuid.NewV4()
		userID, _ := uuid.NewV4()
		userRepo := &mock.MockUserRepo{
			MockUser: []models.User{
				{ID: userID, Name: "Aminat Usman", Email: "meenah20@gmail.com"},
			},
		}
		bookRepo := &mock.MockBookRepository{
			MockBooks: []models.Book{
				{ID: bookID, Title: "Test Book", Author: "Test Author", Status: "Available"},
			},
		}
		service := UserServices{UserRepo: userRepo, BookRepo: bookRepo}

		request := Dto.BookActionRequest{UserID: userID, BookID: bookID, UserName: "Aminat Agbaosi"}
		action, err := service.ReturnBook(request)
		assert.Error(t, err)
		assert.Nil(t, action)
		assert.Equal(t, "User Not Found or Name Doesn't Match", err.Error())
	})

	t.Run("Book Not Avaliable", func(t *testing.T) {
		bookID, _ := uuid.NewV4()
		userID, _ := uuid.NewV4()
		userRepo := &mock.MockUserRepo{
			MockUser: []models.User{
				{ID: userID, Name: "Aminat Usman", Email: "meenah20@gmail.com"},
			},
		}
		bookRepo := &mock.MockBookRepository{
			MockBooks: []models.Book{
				{ID: bookID, Title: "Test Book", Author: "Test Author", Status: "borrowed"},
			},
		}
		service := UserServices{UserRepo: userRepo, BookRepo: bookRepo}
		request := Dto.BookActionRequest{UserID: userID, BookID: bookID, UserName: "Aminat Usman"}
		action, err := service.ReserveBook(request)
		assert.Error(t, err)
		assert.Nil(t, action)
		assert.Equal(t, "Book is not Available", err.Error())

	})

}
