package loan_repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/iwansofian0512/billing_service/internal/model"
	"github.com/jmoiron/sqlx"
)

func newMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	return sqlx.NewDb(db, "postgres"), mock
}

func TestPostgresLoanRepository_UpdateLoan(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	repo := NewPostgresLoanRepository(db)

	loan := &model.Loan{
		ID:                1,
		OutstandingAmount: 0,
		IsActive:          false,
		Status:            model.LoanStatusCompleted,
	}

	query := regexp.QuoteMeta(`UPDATE loans SET outstanding_amount = $1, is_active = $2, status = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4`)

	mock.ExpectExec(query).
		WithArgs(loan.OutstandingAmount, loan.IsActive, loan.Status, loan.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateLoan(context.Background(), loan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestPostgresLoanRepository_UpdateSchedule(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	repo := NewPostgresLoanRepository(db)

	schedule := &model.BillingSchedule{
		ID:         1,
		AmountPaid: 110000,
		Status:     model.BillingStatusPaid,
	}

	query := regexp.QuoteMeta(`UPDATE billing_schedules SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`)

	mock.ExpectExec(query).
		WithArgs(schedule.Status, schedule.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateSchedule(context.Background(), schedule)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestPostgresLoanRepository_GetCurrentPendingSchedules(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	repo := NewPostgresLoanRepository(db)

	query := regexp.QuoteMeta(`WITH pending AS (
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
              ORDER BY week_number ASC`)
	rows := sqlmock.NewRows([]string{"id", "loan_id", "week_number", "due_date", "amount_due", "amount_paid", "status"}).
		AddRow(1, 1, 1, time.Now(), 110000, 0, model.BillingStatusPending).
		AddRow(2, 1, 2, time.Now(), 110000, 0, model.BillingStatusPending)

	mock.ExpectQuery(query).
		WithArgs(1).
		WillReturnRows(rows)

	schedules, err := repo.GetCurrentPendingSchedules(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(schedules) != 2 {
		t.Fatalf("expected 2 schedules, got %d", len(schedules))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}
