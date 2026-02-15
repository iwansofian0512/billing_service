package payment_service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/iwansofian0512/billing_service/internal/model"
)

type mockLoanRepo struct {
	loan              *model.Loan
	schedules         []model.BillingSchedule
	updateLoanErr     error
	updateScheduleErr error
}

func (m *mockLoanRepo) CreateLoan(_ context.Context, loan *model.Loan) error {
	m.loan = loan
	return nil
}

func (m *mockLoanRepo) GetByID(_ context.Context, id int) (*model.Loan, error) {
	return m.loan, nil
}

func (m *mockLoanRepo) GetLoanByID(_ context.Context, id int) (*model.Loan, error) {
	return m.loan, nil
}

func (m *mockLoanRepo) UpdateLoan(_ context.Context, loan *model.Loan) error {
	if m.updateLoanErr != nil {
		return m.updateLoanErr
	}
	m.loan = loan
	return nil
}

func (m *mockLoanRepo) GetSchedules(_ context.Context, loanID int) ([]model.BillingSchedule, error) {
	return m.schedules, nil
}

func (m *mockLoanRepo) GetCurrentPendingSchedules(_ context.Context, loanID int) ([]model.BillingSchedule, error) {
	var result []model.BillingSchedule
	now := time.Now()
	for _, s := range m.schedules {
		if s.Status == model.BillingStatusPending && !s.DueDate.After(now) {
			result = append(result, s)
		}
	}
	return result, nil
}

func (m *mockLoanRepo) GetBorrowerLoans(_ context.Context, borrowerID, page, pageSize int) ([]model.Loan, error) {
	if m.loan != nil {
		return []model.Loan{*m.loan}, nil
	}
	return []model.Loan{}, nil
}

func (m *mockLoanRepo) GetActiveLoanByID(_ context.Context, id int) (*model.Loan, error) {
	return m.loan, nil
}

func (m *mockLoanRepo) UpdateSchedule(_ context.Context, s *model.BillingSchedule) error {
	if m.updateScheduleErr != nil {
		return m.updateScheduleErr
	}
	for i := range m.schedules {
		if m.schedules[i].ID == s.ID {
			m.schedules[i] = *s
			return nil
		}
	}
	return nil
}

type mockPaymentRepo struct {
	lastPayment *model.Payment
	addErr      error
}

func (m *mockPaymentRepo) AddPayment(_ context.Context, p *model.Payment) error {
	if m.addErr != nil {
		return m.addErr
	}
	m.lastPayment = p
	return nil
}

func TestPaymentService_MakePayment(t *testing.T) {
	baseLoan := &model.Loan{
		ID:                  1,
		OutstandingAmount:   5500000,
		WeeklyPaymentAmount: 110000,
		IsActive:            true,
		Status:              model.LoanStatusInProgress,
	}

	loanRepo := &mockLoanRepo{
		loan: &model.Loan{
			ID:                  baseLoan.ID,
			OutstandingAmount:   baseLoan.OutstandingAmount,
			WeeklyPaymentAmount: baseLoan.WeeklyPaymentAmount,
			IsActive:            baseLoan.IsActive,
			Status:              baseLoan.Status,
		},
		schedules: []model.BillingSchedule{
			{ID: 1, WeekNumber: 1, AmountDue: 110000, Status: model.BillingStatusPending, DueDate: time.Now().AddDate(0, 0, -7)},
			{ID: 2, WeekNumber: 2, AmountDue: 110000, Status: model.BillingStatusPending, DueDate: time.Now().AddDate(0, 0, 7)},
		},
	}
	paymentRepo := &mockPaymentRepo{}
	svc := NewPaymentService(loanRepo, paymentRepo)

	t.Run("successful payment", func(t *testing.T) {
		err := svc.MakePayment(context.Background(), 1, 110000)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if loanRepo.loan.OutstandingAmount != 5390000 {
			t.Errorf("expected outstanding 5390000, got %v", loanRepo.loan.OutstandingAmount)
		}

		if loanRepo.schedules[0].Status != model.BillingStatusPaid {
			t.Errorf("expected schedule 1 to be paid")
		}

		if paymentRepo.lastPayment == nil {
			t.Fatalf("expected payment to be recorded")
		}
		if paymentRepo.lastPayment.BillingScheduleID != 1 {
			t.Fatalf("expected billing_schedule_id 1, got %d", paymentRepo.lastPayment.BillingScheduleID)
		}
	})

	t.Run("wrong amount", func(t *testing.T) {
		err := svc.MakePayment(context.Background(), 1, 100000)
		if err == nil {
			t.Errorf("expected error for wrong amount")
		}
	})

	t.Run("late payments require full amount", func(t *testing.T) {
		loanRepo.schedules = []model.BillingSchedule{
			{ID: 1, WeekNumber: 1, AmountDue: 110000, Status: model.BillingStatusPending, DueDate: time.Now().AddDate(0, 0, -14)},
			{ID: 2, WeekNumber: 2, AmountDue: 110000, Status: model.BillingStatusPending, DueDate: time.Now().AddDate(0, 0, -7)},
		}
		err := svc.MakePayment(context.Background(), 1, 110000)
		if err == nil {
			t.Errorf("expected error for partial late payment")
		}

		err = svc.MakePayment(context.Background(), 1, 220000)
		if err != nil {
			t.Fatalf("unexpected error for full late payment: %v", err)
		}
	})

	t.Run("rollback on AddPayment error", func(t *testing.T) {
		loanRepo.loan = &model.Loan{
			ID:                  baseLoan.ID,
			OutstandingAmount:   baseLoan.OutstandingAmount,
			WeeklyPaymentAmount: baseLoan.WeeklyPaymentAmount,
			IsActive:            baseLoan.IsActive,
			Status:              baseLoan.Status,
		}
		loanRepo.schedules = []model.BillingSchedule{
			{ID: 1, WeekNumber: 1, AmountDue: 110000, Status: model.BillingStatusPending, DueDate: time.Now().AddDate(0, 0, -7)},
			{ID: 2, WeekNumber: 2, AmountDue: 110000, Status: model.BillingStatusPending, DueDate: time.Now().AddDate(0, 0, -7)},
		}
		paymentRepo.addErr = errors.New("add payment error")
		paymentRepo.lastPayment = nil

		err := svc.MakePayment(context.Background(), 1, 220000)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		if loanRepo.loan.OutstandingAmount != baseLoan.OutstandingAmount {
			t.Fatalf("expected outstanding %v, got %v", baseLoan.OutstandingAmount, loanRepo.loan.OutstandingAmount)
		}

		for _, s := range loanRepo.schedules {
			if s.Status != model.BillingStatusPending {
				t.Fatalf("expected schedule %d to be pending after rollback, got %s", s.ID, s.Status)
			}
		}

		if paymentRepo.lastPayment != nil {
			t.Fatalf("expected no payment to be recorded on error")
		}

		paymentRepo.addErr = nil
	})

	t.Run("rollback on UpdateLoan error", func(t *testing.T) {
		loanRepo.loan = &model.Loan{
			ID:                  baseLoan.ID,
			OutstandingAmount:   baseLoan.OutstandingAmount,
			WeeklyPaymentAmount: baseLoan.WeeklyPaymentAmount,
			IsActive:            baseLoan.IsActive,
			Status:              baseLoan.Status,
		}
		loanRepo.schedules = []model.BillingSchedule{
			{ID: 1, WeekNumber: 1, AmountDue: 110000, Status: model.BillingStatusPending, DueDate: time.Now().AddDate(0, 0, -7)},
		}
		loanRepo.updateLoanErr = errors.New("update loan error")

		err := svc.MakePayment(context.Background(), 1, 110000)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		if loanRepo.loan.OutstandingAmount != baseLoan.OutstandingAmount {
			t.Fatalf("expected outstanding %v after rollback, got %v", baseLoan.OutstandingAmount, loanRepo.loan.OutstandingAmount)
		}

		if loanRepo.loan.Status != baseLoan.Status {
			t.Fatalf("expected loan status %v after rollback, got %v", baseLoan.Status, loanRepo.loan.Status)
		}

		for _, s := range loanRepo.schedules {
			if s.Status != model.BillingStatusPending {
				t.Fatalf("expected schedule %d to be pending after rollback, got %s", s.ID, s.Status)
			}
		}

		loanRepo.updateLoanErr = nil
	})
}
