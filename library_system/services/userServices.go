package services

import (
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"library-system/Dto"
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
	log.Printf("Starting checkout process for book ID: %v by email: %v", request.BookID, request.Email)

	normalizedEmail := normalizeEmail(request.Email)
	if !isValidEmail(normalizedEmail) {
		return nil, errors.New("Invalid Email Address")
	}

	user, err := s.UserRepo.GetUserByEmail(normalizedEmail)
	if err != nil {
		return nil, fmt.Errorf("User not found: %v", err)
	}

	book, err := s.BookRepo.GetBookByID(request.BookID)
	if err != nil {
		return nil, fmt.Errorf("Book not found: %v", err)
	}

	if strings.ToLower(book.Status) != "available" {
		log.Printf("Book is currently %s", book.Status)
		return nil, fmt.Errorf("Book is not available for checkout")
	}

	existingLoan, err := s.LoanRepo.GetLoanByBookAndEmail(request.BookID, normalizedEmail)
	if err == nil && existingLoan != nil && existingLoan.ReturnDate == nil {
		return nil, errors.New("You have already borrowed this book")
	}

	now := time.Now()
	loan := &models.Loan{
		ID:         uuid.Must(uuid.NewV4()),
		BookID:     request.BookID,
		Email:      normalizedEmail,
		UserID:     user.ID,
		LoanDate:   now,
		ReturnDate: nil,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	book.Status = "borrowed"
	if err := s.BookRepo.UpdateBook(book); err != nil {
		log.Printf("Failed to update book status: %v", err)
		return nil, fmt.Errorf("Failed to update book status: %v", err)
	}

	if err := s.LoanRepo.AddLoan(loan); err != nil {
		book.Status = "available"
		_ = s.BookRepo.UpdateBook(book)
		return nil, fmt.Errorf("Failed to create loan: %v", err)
	}

	return &Dto.BookActionResponse{
		ID:       loan.ID,
		UserID:   user.ID,
		BookID:   loan.BookID,
		Email:    loan.Email,
		Status:   "borrowed",
		LoanDate: loan.LoanDate,
	}, nil
}

func (s *UserServices) ReturnBook(request Dto.BookActionRequest) (*Dto.BookActionResponse, error) {
	log.Printf("Starting return process for book ID: %v by email: %v", request.BookID, request.Email)

	normalizedEmail := normalizeEmail(request.Email)
	if !isValidEmail(normalizedEmail) {
		return nil, errors.New("Invalid Email Address")
	}

	loan, err := s.LoanRepo.GetLoanByBookAndEmail(request.BookID, normalizedEmail)
	if err != nil || loan == nil {
		return nil, errors.New("No active loan found for this book")
	}

	if loan.ReturnDate != nil {
		return nil, errors.New("Book has already been returned")
	}

	book, err := s.BookRepo.GetBookByID(request.BookID)
	if err != nil {
		return nil, fmt.Errorf("Book not found: %v", err)
	}

	now := time.Now()
	loan.ReturnDate = &now
	loan.UpdatedAt = now

	book.Status = "available"
	if err := s.BookRepo.UpdateBook(book); err != nil {
		return nil, fmt.Errorf("Failed to update book status: %v", err)
	}

	if err := s.LoanRepo.UpdateLoan(loan); err != nil {
		book.Status = "borrowed"
		_ = s.BookRepo.UpdateBook(book)
		return nil, fmt.Errorf("Failed to update loan: %v", err)
	}

	return &Dto.BookActionResponse{
		ID:         loan.ID,
		BookID:     loan.BookID,
		Email:      loan.Email,
		Status:     "available",
		LoanDate:   loan.LoanDate,
		ReturnDate: loan.ReturnDate,
	}, nil
}

func (s *UserServices) ReserveBook(request Dto.BookActionRequest) (*Dto.BookActionResponse, error) {
	log.Printf("Starting reservation process for book ID: %v by email: %v", request.BookID, request.Email)

	normalizedEmail := normalizeEmail(request.Email)
	if !isValidEmail(normalizedEmail) {
		return nil, errors.New("Invalid Email Address")
	}

	user, err := s.UserRepo.GetUserByEmail(normalizedEmail)
	if err != nil {
		return nil, fmt.Errorf("User not found: %v", err)
	}

	book, err := s.BookRepo.GetBookByID(request.BookID)
	if err != nil {
		return nil, fmt.Errorf("Book not found: %v", err)
	}

	if strings.ToLower(book.Status) != "available" {
		return nil, errors.New("Book is not available for reservation")
	}

	now := time.Now()
	loan := &models.Loan{
		ID:        uuid.Must(uuid.NewV4()),
		BookID:    request.BookID,
		Email:     normalizedEmail,
		UserID:    user.ID,
		LoanDate:  now,
		CreatedAt: now,
		UpdatedAt: now,
	}

	book.Status = "reserved"
	if err := s.BookRepo.UpdateBook(book); err != nil {
		return nil, fmt.Errorf("Failed to update book status: %v", err)
	}

	if err := s.LoanRepo.AddLoan(loan); err != nil {
		book.Status = "available"
		_ = s.BookRepo.UpdateBook(book)
		return nil, fmt.Errorf("Failed to create reservation: %v", err)
	}

	return &Dto.BookActionResponse{
		ID:       loan.ID,
		UserID:   user.ID,
		BookID:   loan.BookID,
		Email:    loan.Email,
		Status:   "reserved",
		LoanDate: loan.LoanDate,
	}, nil
}

// Helper functions

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
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

//func isStatusAvailable(status string) bool {
//	return normalizeStatus(status) == "available"
//}

func normalizeStatus(status string) string {

	return strings.ToLower(strings.TrimSpace(status))
}
