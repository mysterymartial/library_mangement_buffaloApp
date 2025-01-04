package mock

import (
	"errors"
	"github.com/gofrs/uuid"
	"library-system/models"
)

type MockLoanRepository struct {
	MockLoans                 []models.Loan
	AddLoanError              error
	GetLoanByBookAndUserError error
	UpdateLoanError           error
}

func (r *MockLoanRepository) AddLoan(loan *models.Loan) error {
	if r.AddLoanError != nil {
		return r.AddLoanError
	}
	r.MockLoans = append(r.MockLoans, *loan)
	return nil
}

func (r *MockLoanRepository) GetLoanByBookAndUser(bookID, userID uuid.UUID) (*models.Loan, error) {
	if r.GetLoanByBookAndUserError != nil {
		return nil, r.GetLoanByBookAndUserError
	}
	for _, loan := range r.MockLoans {
		if loan.BookID == bookID && loan.UserID == userID && loan.ReturnDate == nil {
			return &loan, nil
		}
	}
	return nil, errors.New("loan not found")
}

func (r *MockLoanRepository) UpdateLoan(loan *models.Loan) error {
	if r.UpdateLoanError != nil {
		return r.UpdateLoanError
	}
	for i, existingLoan := range r.MockLoans {
		if existingLoan.ID == loan.ID {
			r.MockLoans[i] = *loan
			return nil
		}
	}
	return errors.New("loan not found")
}
