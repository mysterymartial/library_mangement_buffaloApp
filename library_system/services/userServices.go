package services

import (
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"library-system/Dto"
	"library-system/mapper"
	"library-system/models"
	"library-system/repositories/repository"
	"log"
	"regexp"
	"strings"
	"time"
)

type UserServices struct {
	UserRepo repository.UserRepository
	LoanRepo repository.LoanRepository
	BookRepo repository.BookRepository
}

func (s *UserServices) RegisterUser(request Dto.UserRequest) (*Dto.UserResponse, error) {

	normalizedName := normalizeName(request.Name)
	if !isNameValid(normalizedName) {
		return nil, errors.New("Invalid Name")
	}

	normalizedEmail := normalizeEmail(request.Email)

	if !isValidEmail(normalizedEmail) {
		return nil, errors.New("Invalid Email Address")
	}

	existingUser, err := s.UserRepo.GetUserByEmail(normalizedEmail)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	user := &models.User{
		Name:  normalizedName,
		Email: normalizedEmail,
	}

	if err := s.UserRepo.AddUser(user); err != nil {
		log.Printf("Error adding user with email %v: %v", normalizedEmail, err)
		return nil, err
	}

	return &Dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *UserServices) CheckOutBook(request Dto.BookActionRequest) (*Dto.BookActionResponse, error) {
	log.Printf("Starting checkout process for book ID: %v by user ID: %v", request.BookID, request.UserID)

	user, err := s.UserRepo.GetUserByID(request.UserID)
	if err != nil {
		log.Printf("Error fetching user with ID %v: %v", request.UserID, err)
		return nil, errors.New("User not found")
	}

	log.Printf("Fetched user: %v", user)

	if normalizeName(user.Name) != normalizeName(request.UserName) {
		log.Printf("User name mismatch: %v != %v", user.Name, request.UserName)
		return nil, errors.New("Invalid User Name")
	}

	book, err := s.BookRepo.GetBookByID(request.BookID)
	if err != nil {
		log.Printf("Error fetching book with ID %v: %v", request.BookID, err)
		return nil, errors.New("Book not found")
	}

	log.Printf("Fetched book: %v", book)

	existingLoan, err := s.LoanRepo.GetLoanByBookAndUser(request.BookID, request.UserID)
	if err == nil && existingLoan != nil && existingLoan.ReturnDate == nil {
		return nil, errors.New("You have already borrowed this book")
	}

	if book.Status != "available" {
		log.Printf("Book with ID %v is not available. Current status: %v", request.BookID, book.Status)
		return nil, fmt.Errorf("Book is currently %s and cannot be checked out", book.Status)
	}

	now := time.Now()
	loan := &models.Loan{
		ID:         uuid.Must(uuid.NewV4()),
		BookID:     request.BookID,
		UserID:     request.UserID,
		UserName:   request.UserName,
		LoanDate:   now,
		ReturnDate: nil,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	book.Status = "borrowed"
	book.UpdatedAt = now
	if err := s.BookRepo.UpdateBook(book); err != nil {
		log.Printf("Error updating book status for book ID %v: %v", request.BookID, err)
		return nil, fmt.Errorf("Failed to update book status: %v", err)
	}

	if err := s.LoanRepo.AddLoan(loan); err != nil {

		book.Status = "available"
		_ = s.BookRepo.UpdateBook(book)
		log.Printf("Error adding loan record for book ID %v: %v", request.BookID, err)
		return nil, fmt.Errorf("Failed to create loan record: %v", err)
	}

	log.Printf("Successfully created loan: %v", loan)
	return mapper.ToBookActionResponse(loan, "borrowed"), nil
}

func (s *UserServices) ReserveBook(request Dto.BookActionRequest) (*Dto.BookActionResponse, error) {
	log.Printf("Starting reservation process for book ID: %v by user ID: %v", request.BookID, request.UserID)
	// Fetch the user and validate using normalized name
	user, err := s.UserRepo.GetUserByID(request.UserID)
	if err != nil || normalizeName(user.Name) != normalizeName(request.UserName) {
		log.Printf("User validation failed for user ID %v: %v", request.UserID, err)
		return nil, errors.New("User Not Found or Name Doesn't Match")
	}

	book, err := s.BookRepo.GetBookByID(request.BookID)
	if err != nil {
		log.Printf("Error fetching book with ID %v: %v", request.BookID, err)
		return nil, err
	}

	if !isStatusAvailable(book.Status) {
		log.Printf("Book with ID %v is not available. Status: %v", request.BookID, book.Status)
		return nil, errors.New("Book is not Available")
	}

	book.Status = "Reserved"
	if err := s.BookRepo.UpdateBook(book); err != nil {
		log.Printf("Error updating book status for book ID %v: %v", request.BookID, err)
		return nil, err
	}

	return mapper.ToBookActionResponse(&models.Loan{
		BookID:   request.BookID,
		UserID:   request.UserID,
		UserName: request.UserName,
	}, "Reserved"), nil
}
func (s *UserServices) ReturnBook(request Dto.BookActionRequest) (*Dto.BookActionResponse, error) {
	log.Printf("Starting return process for book ID: %v by user ID: %v", request.BookID, request.UserID)

	user, err := s.UserRepo.GetUserByID(request.UserID)
	if err != nil || normalizeName(user.Name) != normalizeName(request.UserName) {
		log.Printf("User validation failed for user ID %v: %v", request.UserID, err)
		return nil, errors.New("User Not Found or Name Doesn't Match")
	}

	book, err := s.BookRepo.GetBookByID(request.BookID)
	if err != nil {
		log.Printf("Error fetching book with ID %v: %v", request.BookID, err)
		return nil, err
	}

	loan, err := s.LoanRepo.GetLoanByBookAndUser(request.BookID, request.UserID)
	if err != nil {
		log.Printf("Loan record not found for book ID %v and user ID %v: %v", request.BookID, request.UserID, err)
		return nil, errors.New("Loan Record Not Found")
	}

	if loan.ReturnDate != nil {
		log.Printf("Book with ID %v has already been returned", request.BookID)
		return nil, errors.New("Book has already been returned")
	}

	loan.ReturnDate = new(time.Time)
	*loan.ReturnDate = time.Now()

	if err := s.LoanRepo.UpdateLoan(loan); err != nil {
		log.Printf("Error updating loan record for book ID %v: %v", request.BookID, err)
		return nil, err
	}

	// Update the book status to "available"
	book.Status = "available"
	if err := s.BookRepo.UpdateBook(book); err != nil {
		log.Printf("Error updating book status for book ID %v: %v", request.BookID, err)
		return nil, err
	}

	return mapper.ToBookActionResponse(loan, "Returned"), nil
}

func isValidEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(emailRegex)
	return regex.MatchString(email)
}

func isNameValid(name string) bool {
	const nameRegex = `^[a-zA-Z\s]+$`
	regex := regexp.MustCompile(nameRegex)
	return regex.MatchString(name)
}

func normalizeName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func isStatusAvailable(status string) bool {
	return normalizeStatus(status) == "available"
}

func normalizeStatus(status string) string {
	// Trim spaces and handle potential unexpected characters
	return strings.ToLower(strings.TrimSpace(status))
}
