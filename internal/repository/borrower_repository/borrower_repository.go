package borrower_repository

import (
	"context"
	"database/sql"

	"github.com/iwansofian0512/billing_service/internal/model"
	"github.com/jmoiron/sqlx"
)

type postgresBorrowerRepository struct {
	db *sqlx.DB
}

func NewPostgresBorrowerRepository(db *sqlx.DB) BorrowerRepository {
	return &postgresBorrowerRepository{db: db}
}

type BorrowerRepository interface {
	Create(ctx context.Context, b *model.Borrower) error
	GetByEmail(ctx context.Context, email string) (*model.Borrower, error)
}

func (r *postgresBorrowerRepository) Create(ctx context.Context, borrower *model.Borrower) error {
	query := `INSERT INTO borrowers (name, email, is_active) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query, borrower.Name, borrower.Email, borrower.IsActive).Scan(&borrower.ID, &borrower.CreatedAt, &borrower.UpdatedAt)
}

func (r *postgresBorrowerRepository) GetByEmail(ctx context.Context, email string) (*model.Borrower, error) {
	var b model.Borrower
	query := `SELECT id, name, email, is_active, created_at, updated_at FROM borrowers WHERE email = $1`
	err := r.db.GetContext(ctx, &b, query, email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &b, nil
}
