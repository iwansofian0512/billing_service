package borrower_repository

import (
	"context"
	"database/sql"
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

func TestPostgresBorrowerRepository_Create(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	repo := NewPostgresBorrowerRepository(db)

	b := &model.Borrower{
		Name:     "John Doe",
		Email:    "john@example.com",
		IsActive: true,
	}

	query := regexp.QuoteMeta(`INSERT INTO borrowers (name, email, is_active) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`)
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(1, time.Now(), time.Now())

	mock.ExpectQuery(query).
		WithArgs(b.Name, b.Email, b.IsActive).
		WillReturnRows(rows)

	err := repo.Create(context.Background(), b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if b.ID != 1 {
		t.Fatalf("expected id 1, got %d", b.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestPostgresBorrowerRepository_GetByEmail_Found(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	repo := NewPostgresBorrowerRepository(db)

	query := regexp.QuoteMeta(`SELECT id, name, email, is_active, created_at, updated_at FROM borrowers WHERE email = $1`)
	rows := sqlmock.NewRows([]string{"id", "name", "email", "is_active"}).
		AddRow(1, "John Doe", "john@example.com", true)

	mock.ExpectQuery(query).
		WithArgs("john@example.com").
		WillReturnRows(rows)

	b, err := repo.GetByEmail(context.Background(), "john@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if b == nil || b.Email != "john@example.com" {
		t.Fatalf("unexpected borrower: %+v", b)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestPostgresBorrowerRepository_GetByEmail_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	repo := NewPostgresBorrowerRepository(db)

	query := regexp.QuoteMeta(`SELECT id, name, email, is_active, created_at, updated_at FROM borrowers WHERE email = $1`)

	mock.ExpectQuery(query).
		WithArgs("missing@example.com").
		WillReturnError(sql.ErrNoRows)

	b, err := repo.GetByEmail(context.Background(), "missing@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if b != nil {
		t.Fatalf("expected nil borrower, got %+v", b)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}
