package loan_repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/iwansofian0512/billing_service/internal/model"
	"github.com/jmoiron/sqlx"
)

type postgresLoanRepository struct {
	db *sqlx.DB
}

func NewPostgresLoanRepository(db *sqlx.DB) LoanRepository {
	return &postgresLoanRepository{db: db}
}

type LoanRepository interface {
	CreateLoan(ctx context.Context, loan *model.Loan) error
	GetActiveLoanByID(ctx context.Context, id int) (*model.Loan, error)
	UpdateLoan(ctx context.Context, loan *model.Loan) error
	GetCurrentPendingSchedules(ctx context.Context, loanID int) ([]model.BillingSchedule, error)
	GetBorrowerLoans(ctx context.Context, borrowerID, page, pageSize int) ([]model.Loan, error)
	UpdateSchedule(ctx context.Context, schedule *model.BillingSchedule) error
}

func (r *postgresLoanRepository) CreateLoan(ctx context.Context, loan *model.Loan) error {
	query := `INSERT INTO loans (borrower_id, principal_amount, total_interest, total_payable, outstanding_amount, duration_weeks, weekly_payment_amount, is_active, status)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, created_at, updated_at`
	err := r.db.QueryRowContext(ctx, query, loan.BorrowerID, loan.PrincipalAmount, loan.TotalInterest, loan.TotalPayable, loan.OutstandingAmount, loan.DurationWeeks, loan.WeeklyPaymentAmount, loan.IsActive, loan.Status).
		Scan(&loan.ID, &loan.CreatedAt, &loan.UpdatedAt)
	if err != nil {
		return err
	}

	for _, s := range loan.Schedules {
		queryS := `INSERT INTO billing_schedules (loan_id, week_number, due_date, amount_due, amount_paid, status)
                   VALUES ($1, $2, $3, $4, $5, $6)`
		_, err = r.db.ExecContext(ctx, queryS, loan.ID, s.WeekNumber, s.DueDate, s.AmountDue, s.AmountPaid, s.Status)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *postgresLoanRepository) GetActiveLoanByID(ctx context.Context, id int) (*model.Loan, error) {
	var loan model.Loan
	query := `SELECT id, borrower_id, principal_amount, total_interest, total_payable, outstanding_amount, duration_weeks, weekly_payment_amount, is_active, status, created_at, updated_at
              FROM loans WHERE id = $1 AND is_active = TRUE AND status = 'inprogress'`
	err := r.db.GetContext(ctx, &loan, query, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("active loan not found")
	}
	if err != nil {
		return nil, err
	}
	return &loan, nil
}

func (r *postgresLoanRepository) GetBorrowerLoans(ctx context.Context, borrowerID, page, pageSize int) ([]model.Loan, error) {
	var loans []model.Loan
	if pageSize <= 0 {
		pageSize = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * pageSize

	query := `SELECT
                l.id,
                l.borrower_id,
                l.principal_amount,
                l.total_interest,
                l.total_payable,
                l.outstanding_amount,
                l.duration_weeks,
                l.weekly_payment_amount,
                l.is_active,
                l.status,
                l.created_at,
                l.updated_at,
                CASE
                    WHEN (
                        SELECT COUNT(*)
                        FROM billing_schedules bs
                        WHERE bs.loan_id = l.id
                          AND bs.status = 'pending'
                          AND bs.due_date < CURRENT_DATE
                    ) >= 2 THEN TRUE
                    ELSE FALSE
                END AS is_delinquent
            FROM loans l`

	var err error
	if borrowerID > 0 {
		query += ` WHERE l.borrower_id = $1 ORDER BY l.created_at DESC LIMIT $2 OFFSET $3`
		err = r.db.SelectContext(ctx, &loans, query, borrowerID, pageSize, offset)
	} else {
		query += ` ORDER BY l.created_at DESC LIMIT $1 OFFSET $2`
		err = r.db.SelectContext(ctx, &loans, query, pageSize, offset)
	}
	return loans, err
}

func (r *postgresLoanRepository) UpdateLoan(ctx context.Context, loan *model.Loan) error {
	query := `UPDATE loans SET outstanding_amount = $1, is_active = $2, status = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4`
	_, err := r.db.ExecContext(ctx, query, loan.OutstandingAmount, loan.IsActive, loan.Status, loan.ID)
	return err
}

func (r *postgresLoanRepository) GetCurrentPendingSchedules(ctx context.Context, loanID int) ([]model.BillingSchedule, error) {
	var schedules []model.BillingSchedule
	query := `WITH pending AS (
                SELECT id, loan_id, week_number, due_date, amount_due, amount_paid, status, created_at, updated_at
                FROM billing_schedules
                WHERE loan_id = $1 AND status = 'pending'
              ),
              overdue AS (
                SELECT *
                FROM pending
                WHERE due_date <= CURRENT_DATE
              ),
              next_upcoming AS (
                SELECT *
                FROM pending
                WHERE due_date > CURRENT_DATE
                ORDER BY due_date ASC
                LIMIT 1
              )
              SELECT id, loan_id, week_number, due_date, amount_due, amount_paid, status, created_at, updated_at
              FROM overdue
              UNION ALL
              SELECT id, loan_id, week_number, due_date, amount_due, amount_paid, status, created_at, updated_at
              FROM next_upcoming
              ORDER BY week_number ASC`
	err := r.db.SelectContext(ctx, &schedules, query, loanID)
	return schedules, err
}

func (r *postgresLoanRepository) UpdateSchedule(ctx context.Context, s *model.BillingSchedule) error {
	query := `UPDATE billing_schedules SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, s.Status, s.ID)
	return err
}
