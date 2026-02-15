package borrower_service

import (
	"context"
	"errors"

	"github.com/iwansofian0512/billing_service/internal/model"
	"github.com/iwansofian0512/billing_service/internal/repository/borrower_repository"
	"github.com/iwansofian0512/billing_service/internal/repository/loan_repository"
	"github.com/iwansofian0512/billing_service/internal/service/loan_service"
)

type BorrowerService interface {
	CreateBorrower(ctx context.Context, name, email string) (*model.Borrower, error)
	ListBorrowerLoans(ctx context.Context, borrowerID, page, pageSize int) ([]model.Loan, error)
}

type borrowerService struct {
	borrowerRepo borrower_repository.BorrowerRepository
	loanRepo     loan_repository.LoanRepository
	loanService  loan_service.LoanService
}

var ErrBorrowerEmailExists = errors.New("email already registered")

func NewBorrowerService(borrowerRepo borrower_repository.BorrowerRepository, loanRepo loan_repository.LoanRepository, loanService loan_service.LoanService) BorrowerService {
	return &borrowerService{
		borrowerRepo: borrowerRepo,
		loanRepo:     loanRepo,
		loanService:  loanService,
	}
}

func (s *borrowerService) CreateBorrower(ctx context.Context, name, email string) (*model.Borrower, error) {
	existing, err := s.borrowerRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrBorrowerEmailExists
	}

	borrower := &model.Borrower{
		Name:  name,
		Email: email,
	}

	if err := s.borrowerRepo.Create(ctx, borrower); err != nil {
		return nil, err
	}

	return borrower, nil
}

func (s *borrowerService) ListBorrowerLoans(ctx context.Context, borrowerID, page, pageSize int) ([]model.Loan, error) {
	if pageSize <= 0 {
		pageSize = 10
	}
	if page <= 0 {
		page = 1
	}

	loans, err := s.loanRepo.GetBorrowerLoans(ctx, borrowerID, page, pageSize)
	if err != nil {
		return nil, err
	}

	return loans, nil
}
