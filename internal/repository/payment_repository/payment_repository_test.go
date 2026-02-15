package payment_repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/iwansofian0512/billing_service/internal/model"
)

func newMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	return sqlx.NewDb(db, "postgres"), mock
}

func TestPostgresPaymentRepository_AddPayment(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	repo := NewPostgresPaymentRepository(db)

	p := &model.Payment{
		LoanID:            1,
		BillingScheduleID: 10,
		Amount:            110000,
		PaymentDate:       time.Now(),
	}

	query := regexp.QuoteMeta(`INSERT INTO payments (loan_id, billing_schedule_id, amount, payment_date) VALUES ($1, $2, $3, $4) RETURNING id`)

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(1)

	mock.ExpectQuery(query).
		WithArgs(p.LoanID, p.BillingScheduleID, p.Amount, p.PaymentDate).
		WillReturnRows(rows)

	err := repo.AddPayment(context.Background(), p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p.ID != 1 {
		t.Fatalf("expected id 1, got %d", p.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}
