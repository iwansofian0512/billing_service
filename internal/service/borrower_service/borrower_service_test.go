package borrower_service

import (
	"context"
	"errors"
	"testing"

	"github.com/iwansofian0512/billing_service/internal/model"
)

type mockBorrowerRepo struct {
	created   *model.Borrower
	shouldErr bool
	existing  *model.Borrower
}

func (m *mockBorrowerRepo) Create(_ context.Context, b *model.Borrower) error {
	if m.shouldErr {
		return errors.New("create borrower error")
	}
	m.created = b
	return nil
}

func (m *mockBorrowerRepo) GetByID(id int) (*model.Borrower, error) {
	return nil, nil
}

func (m *mockBorrowerRepo) GetByEmail(_ context.Context, email string) (*model.Borrower, error) {
	return m.existing, nil
}

type mockLoanRepo struct {
	loans []model.Loan
}

func (m *mockLoanRepo) CreateLoan(ctx context.Context, loan *model.Loan) error {
	return nil
}

func (m *mockLoanRepo) GetLoanByID(ctx context.Context, id int) (*model.Loan, error) {
	return nil, nil
}

func (m *mockLoanRepo) GetActiveLoanByID(ctx context.Context, id int) (*model.Loan, error) {
	return nil, nil
}

func (m *mockLoanRepo) UpdateLoan(ctx context.Context, loan *model.Loan) error {
	return nil
}

func (m *mockLoanRepo) GetSchedules(ctx context.Context, loanID int) ([]model.BillingSchedule, error) {
	return nil, nil
}

func (m *mockLoanRepo) GetCurrentPendingSchedules(ctx context.Context, loanID int) ([]model.BillingSchedule, error) {
	return nil, nil
}

func (m *mockLoanRepo) GetBorrowerLoans(ctx context.Context, borrowerID, page, pageSize int) ([]model.Loan, error) {
	return m.loans, nil
}

func (m *mockLoanRepo) UpdateSchedule(ctx context.Context, s *model.BillingSchedule) error {
	return nil
}

type mockLoanService struct {
	isDelinquent bool
	err          error
}

func (m *mockLoanService) CreateLoan(ctx context.Context, borrowerID int, amount float64) (*model.Loan, error) {
	return nil, nil
}

func (m *mockLoanService) GetOutstanding(ctx context.Context, loanID int) (float64, error) {
	return 0, nil
}

func (m *mockLoanService) IsDelinquent(ctx context.Context, loanID int) (bool, error) {
	if m.err != nil {
		return false, m.err
	}
	return m.isDelinquent, nil
}

// mockLoanRepo already satisfies GetActiveLoanByID via method above

func TestBorrowerService_CreateBorrower_Success(t *testing.T) {
	borrowerRepo := &mockBorrowerRepo{}
	loanRepo := &mockLoanRepo{}
	loanSvc := &mockLoanService{}
	svc := NewBorrowerService(borrowerRepo, loanRepo, loanSvc)

	b, err := svc.CreateBorrower(context.Background(), "John Doe", "john@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if b.Name != "John Doe" || b.Email != "john@example.com" {
		t.Fatalf("unexpected borrower data: %+v", b)
	}

	if borrowerRepo.created == nil {
		t.Fatalf("expected borrower to be saved in repository")
	}
}

func TestBorrowerService_CreateBorrower_Error(t *testing.T) {
	borrowerRepo := &mockBorrowerRepo{shouldErr: true}
	loanRepo := &mockLoanRepo{}
	loanSvc := &mockLoanService{}
	svc := NewBorrowerService(borrowerRepo, loanRepo, loanSvc)

	b, err := svc.CreateBorrower(context.Background(), "John Doe", "john@example.com")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if b != nil {
		t.Fatalf("expected nil borrower on error")
	}
}

func TestBorrowerService_CreateBorrower_DuplicateEmail(t *testing.T) {
	borrowerRepo := &mockBorrowerRepo{
		existing: &model.Borrower{
			ID:    1,
			Name:  "Existing",
			Email: "john@example.com",
		},
	}
	loanRepo := &mockLoanRepo{}
	loanSvc := &mockLoanService{}
	svc := NewBorrowerService(borrowerRepo, loanRepo, loanSvc)

	b, err := svc.CreateBorrower(context.Background(), "John Doe", "john@example.com")
	if err == nil {
		t.Fatalf("expected error for duplicate email, got nil")
	}
	if !errors.Is(err, ErrBorrowerEmailExists) {
		t.Fatalf("expected ErrBorrowerEmailExists, got %v", err)
	}
	if b != nil {
		t.Fatalf("expected nil borrower on duplicate email")
	}
}

func TestBorrowerService_ListBorrowerLoans_WithBorrowerID(t *testing.T) {
	borrowerRepo := &mockBorrowerRepo{}
	loanRepo := &mockLoanRepo{
		loans: []model.Loan{
			{ID: 1, BorrowerID: 1, IsDelinquent: true},
			{ID: 2, BorrowerID: 1, IsDelinquent: false},
		},
	}
	loanSvc := &mockLoanService{}
	svc := NewBorrowerService(borrowerRepo, loanRepo, loanSvc)

	loans, err := svc.ListBorrowerLoans(context.Background(), 1, 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(loans) != 2 {
		t.Fatalf("expected 2 loans, got %d", len(loans))
	}

	if !loans[0].IsDelinquent {
		t.Fatalf("expected first loan to be delinquent")
	}

	if loans[1].IsDelinquent {
		t.Fatalf("expected second loan to not be delinquent")
	}
}
