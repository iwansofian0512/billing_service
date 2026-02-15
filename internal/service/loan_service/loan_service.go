package loan_service

import (
	"context"
	"time"

	"github.com/iwansofian0512/billing_service/internal/constant"
	"github.com/iwansofian0512/billing_service/internal/model"
	"github.com/iwansofian0512/billing_service/internal/repository/loan_repository"
)

type loanService struct {
	repo loan_repository.LoanRepository
}

func NewLoanService(repo loan_repository.LoanRepository) LoanService {
	return &loanService{
		repo: repo,
	}
}

type LoanService interface {
	CreateLoan(ctx context.Context, borrowerID int, amount float64) (*model.Loan, error)
}

func (s *loanService) CreateLoan(ctx context.Context, borrowerID int, principal float64) (*model.Loan, error) {
	// ... existing logic ...
	interest := principal * constant.LoanInterest
	totalPayable := principal + interest
	weeklyPayment := totalPayable / constant.MaxLoanDuration
	loan := &model.Loan{
		BorrowerID:          borrowerID,
		PrincipalAmount:     principal,
		TotalInterest:       interest,
		TotalPayable:        totalPayable,
		OutstandingAmount:   totalPayable,
		DurationWeeks:       constant.MaxLoanDuration,
		WeeklyPaymentAmount: weeklyPayment,
		IsActive:            true,
		Status:              model.LoanStatusInProgress,
	}

	now := time.Now()
	for durration := 1; durration <= constant.MaxLoanDuration; durration++ {
		schedule := model.BillingSchedule{
			WeekNumber: durration,
			DueDate:    now.AddDate(0, 0, durration*7),
			AmountDue:  weeklyPayment,
			AmountPaid: 0,
			Status:     model.BillingStatusPending,
		}
		loan.Schedules = append(loan.Schedules, schedule)
	}

	err := s.repo.CreateLoan(ctx, loan)
	if err != nil {
		return nil, err
	}

	return loan, nil
}
