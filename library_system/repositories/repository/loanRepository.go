package repository

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"library-system/models"
)

type LoanRepository interface {
	AddLoan(loan *models.Loan) error
	GetLoanByBookAndUser(bookID, userID uuid.UUID) (*models.Loan, error)
	UpdateLoan(loan *models.Loan) error
	GetLoanByBookAndEmail(bookID uuid.UUID, email string) (*models.Loan, error)
}

type loanRepositoryImpl struct {
	DB *pop.Connection
}

func NewLoanRepository(db *pop.Connection) LoanRepository {
	return &loanRepositoryImpl{DB: db}
}

func (r *loanRepositoryImpl) AddLoan(loan *models.Loan) error {
	return r.DB.Create(loan)
}

func (r *loanRepositoryImpl) UpdateLoan(loan *models.Loan) error {
	return r.DB.Update(loan)
}

func (r *loanRepositoryImpl) GetLoanByBookAndUser(bookID, userID uuid.UUID) (*models.Loan, error) {
	loan := &models.Loan{}
	err := r.DB.Where("book_id = ? AND user_id = ? AND return_date IS NULL", bookID, userID).First(loan)
	if err != nil {
		return nil, err
	}
	return loan, nil
}

func (r *loanRepositoryImpl) GetLoanByBookAndEmail(bookID uuid.UUID, email string) (*models.Loan, error) {
	loan := &models.Loan{}
	query := r.DB.Q()
	err := query.
		Join("users u", "loans.user_id = u.id").
		Where("loans.book_id = ? AND u.email = ? AND loans.return_date IS NULL", bookID, email).
		First(loan)
	if err != nil {
		return nil, err
	}
	return loan, nil
}
