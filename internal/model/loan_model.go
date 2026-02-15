package model

import (
	"time"
)

type LoanStatus string

const (
	LoanStatusInProgress LoanStatus = "inprogress"
	LoanStatusCompleted  LoanStatus = "completed"
)

type CreateLoanRequest struct {
	BorrowerID float64 `json:"borrower_id"`
	Amount     float64 `json:"amount"`
}

type Loan struct {
	ID                  int               `json:"id" db:"id"`
	BorrowerID          int               `json:"borrowerID" db:"borrower_id"`
	PrincipalAmount     float64           `json:"principalAmount" db:"principal_amount"`
	TotalInterest       float64           `json:"totalInterest" db:"total_interest"`
	TotalPayable        float64           `json:"totalPayable" db:"total_payable"`
	OutstandingAmount   float64           `json:"outstandingAmount" db:"outstanding_amount"`
	DurationWeeks       int               `json:"durationWeeks" db:"duration_weeks"`
	WeeklyPaymentAmount float64           `json:"weeklyPaymentAmount" db:"weekly_payment_amount"`
	IsActive            bool              `json:"isActive" db:"is_active"`
	Status              LoanStatus        `json:"status" db:"status"`
	CreatedAt           time.Time         `json:"createdAt" db:"created_at"`
	UpdatedAt           time.Time         `json:"updatedAt" db:"updated_at"`
	IsDelinquent        bool              `json:"isDelinquent,omitempty" db:"is_delinquent"`
	Schedules           []BillingSchedule `json:"schedules,omitempty"`
}

type BillingStatus string

const (
	BillingStatusPending BillingStatus = "pending"
	BillingStatusPaid    BillingStatus = "paid"
)

type BillingSchedule struct {
	ID         int           `json:"id" db:"id"`
	LoanID     int           `json:"loanID" db:"loan_id"`
	WeekNumber int           `json:"weekNumber" db:"week_number"`
	DueDate    time.Time     `json:"dueDate" db:"due_date"`
	AmountDue  float64       `json:"amountDue" db:"amount_due"`
	AmountPaid float64       `json:"amountPaid" db:"amount_paid"`
	Status     BillingStatus `json:"status" db:"status"`
	CreatedAt  time.Time     `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time     `json:"updatedAt" db:"updated_at"`
}

type Payment struct {
	ID                int       `json:"id" db:"id"`
	LoanID            int       `json:"loanID" db:"loan_id"`
	BillingScheduleID int       `json:"billingScheduleID" db:"billing_schedule_id"`
	Amount            float64   `json:"amount" db:"amount"`
	PaymentDate       time.Time `json:"paymentDate" db:"payment_date"`
}
