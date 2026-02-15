package payment_service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/iwansofian0512/billing_service/internal/model"
	"github.com/iwansofian0512/billing_service/internal/repository/loan_repository"
	"github.com/iwansofian0512/billing_service/internal/repository/payment_repository"
)

type paymentService struct {
	loanRepo    loan_repository.LoanRepository
	paymentRepo payment_repository.PaymentRepository

	mu           sync.Mutex
	paymentLocks map[int]*sync.Mutex
}

func NewPaymentService(loanRepo loan_repository.LoanRepository, paymentRepo payment_repository.PaymentRepository) PaymentService {
	return &paymentService{
		loanRepo:     loanRepo,
		paymentRepo:  paymentRepo,
		paymentLocks: make(map[int]*sync.Mutex),
	}
}

type PaymentService interface {
	MakePayment(ctx context.Context, loanID int, amount float64) error
}

func (s *paymentService) MakePayment(ctx context.Context, loanID int, amount float64) error {
	lock := s.getPaymentLock(loanID)
	lock.Lock()
	defer lock.Unlock()

	loan, err := s.loanRepo.GetActiveLoanByID(ctx, loanID)
	if err != nil {
		return err
	}

	schedules, err := s.loanRepo.GetCurrentPendingSchedules(ctx, loanID)
	if err != nil {
		return err
	}

	if len(schedules) == 0 {
		return fmt.Errorf("no pending payments found")
	}

	var totalDue float64
	for _, schedule := range schedules {
		totalDue += schedule.AmountDue
	}

	// mimimum payment for pending payment (if pending payment is 3, then minimum payment is 2 week payment, because the third payment is not on due)
	minimumTotalDue := totalDue - loan.WeeklyPaymentAmount
	if len(schedules) > 1 && totalDue == minimumTotalDue {
		schedules = schedules[:len(schedules)-1]
		totalDue = minimumTotalDue
	}

	// validatre price amount
	if len(schedules) > 1 && amount != totalDue {
		return fmt.Errorf("payment must be exactly %v to cover late payments", totalDue)
	}
	if len(schedules) <= 1 && amount != loan.WeeklyPaymentAmount {
		return fmt.Errorf("payment must be exactly %v", loan.WeeklyPaymentAmount)
	}

	originalLoan := *loan
	updatedSchedules := make([]model.BillingSchedule, 0, len(schedules))

	rollback := func() {
		for _, sOriginal := range updatedSchedules {
			_ = s.loanRepo.UpdateSchedule(ctx, &model.BillingSchedule{
				ID:        sOriginal.ID,
				Status:    sOriginal.Status,
				AmountDue: sOriginal.AmountDue,
			})
		}
		*loan = originalLoan
	}

	for _, schedule := range schedules {
		err = s.loanRepo.UpdateSchedule(ctx, &model.BillingSchedule{
			ID:        schedule.ID,
			Status:    model.BillingStatusPaid,
			AmountDue: schedule.AmountDue,
		})
		if err != nil {
			rollback()
			return err
		}

		updatedSchedules = append(updatedSchedules, schedule)

		err = s.paymentRepo.AddPayment(ctx, &model.Payment{
			LoanID:            loanID,
			BillingScheduleID: schedule.ID,
			Amount:            schedule.AmountDue,
			PaymentDate:       time.Now(),
		})
		if err != nil {
			rollback()
			return err
		}
	}

	// update remaining loan amount
	loan.OutstandingAmount -= amount
	if loan.OutstandingAmount <= 0 {
		loan.OutstandingAmount = 0
		loan.IsActive = false
		loan.Status = model.LoanStatusCompleted
	}

	err = s.loanRepo.UpdateLoan(ctx, loan)
	if err != nil {
		rollback()
		return err
	}

	return nil
}

// prevent loan payment race condition
func (s *paymentService) getPaymentLock(loanID int) *sync.Mutex {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.paymentLocks == nil {
		s.paymentLocks = make(map[int]*sync.Mutex)
	}

	if lock, ok := s.paymentLocks[loanID]; ok {
		return lock
	}

	lock := &sync.Mutex{}
	s.paymentLocks[loanID] = lock
	return lock
}
