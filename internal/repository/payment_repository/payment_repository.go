package payment_repository

import (
	"context"

	"github.com/iwansofian0512/billing_service/internal/model"
	"github.com/jmoiron/sqlx"
)

type postgresPaymentRepository struct {
	db *sqlx.DB
}

func NewPostgresPaymentRepository(db *sqlx.DB) PaymentRepository {
	return &postgresPaymentRepository{db: db}
}

type PaymentRepository interface {
	AddPayment(ctx context.Context, payment *model.Payment) error
}

func (r *postgresPaymentRepository) AddPayment(ctx context.Context, p *model.Payment) error {
	query := `INSERT INTO payments (loan_id, billing_schedule_id, amount, payment_date) VALUES ($1, $2, $3, $4) RETURNING id`
	return r.db.QueryRowxContext(ctx, query, p.LoanID, p.BillingScheduleID, p.Amount, p.PaymentDate).Scan(&p.ID)
}
