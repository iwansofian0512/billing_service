package loan_service

import (
	"context"
	"testing"
	"time"

	"github.com/iwansofian0512/billing_service/internal/model"
)

type mockRepo struct {
	loan      *model.Loan
	schedules []model.BillingSchedule
}

func (m *mockRepo) CreateLoan(_ context.Context, loan *model.Loan) error {
	m.loan = loan
	return nil
}

func (m *mockRepo) GetLoanByID(_ context.Context, id int) (*model.Loan, error) {
	return m.loan, nil
}

func (m *mockRepo) UpdateLoan(_ context.Context, loan *model.Loan) error {
	m.loan = loan
	return nil
}

func (m *mockRepo) GetSchedules(_ context.Context, loanID int) ([]model.BillingSchedule, error) {
	return m.schedules, nil
}

func (m *mockRepo) GetCurrentPendingSchedules(_ context.Context, loanID int) ([]model.BillingSchedule, error) {
	var result []model.BillingSchedule
	now := time.Now()
	for _, s := range m.schedules {
		if s.Status == model.BillingStatusPending && !s.DueDate.After(now) {
			result = append(result, s)
		}
	}
	return result, nil
}

func (m *mockRepo) GetBorrowerLoans(_ context.Context, borrowerID, page, pageSize int) ([]model.Loan, error) {
	if m.loan != nil {
		return []model.Loan{*m.loan}, nil
	}
	return []model.Loan{}, nil
}

func (m *mockRepo) UpdateSchedule(_ context.Context, s *model.BillingSchedule) error {
	for i := range m.schedules {
		if m.schedules[i].ID == s.ID {
			m.schedules[i] = *s
			return nil
		}
	}
	return nil
}

func (m *mockRepo) GetActiveLoanByID(_ context.Context, id int) (*model.Loan, error) {
	return m.loan, nil
}

func TestLoanService_CreateLoan(t *testing.T) {
	repo := &mockRepo{}
	svc := NewLoanService(repo)

	loan, err := svc.CreateLoan(context.Background(), 1, 5000000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if loan.BorrowerID != 1 {
		t.Fatalf("expected borrower id 1, got %d", loan.BorrowerID)
	}

	if loan.PrincipalAmount != 5000000 {
		t.Fatalf("expected principal 5000000, got %v", loan.PrincipalAmount)
	}

	if len(loan.Schedules) != 50 {
		t.Fatalf("expected 50 schedules, got %d", len(loan.Schedules))
	}
}

func TestLoanService_MakePayment(t *testing.T) {}
