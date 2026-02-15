package loan_handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/iwansofian0512/billing_service/internal/model"
	"github.com/iwansofian0512/billing_service/internal/service/loan_service"
)

type mockLoanService struct {
	createResult *model.Loan
	createErr    error
}

func (m *mockLoanService) CreateLoan(ctx context.Context, borrowerID int, amount float64) (*model.Loan, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	if m.createResult != nil {
		return m.createResult, nil
	}
	return &model.Loan{
		ID:         1,
		BorrowerID: borrowerID,
	}, nil
}

func (m *mockLoanService) GetOutstanding(ctx context.Context, loanID int) (float64, error) {
	return 0, nil
}

func (m *mockLoanService) IsDelinquent(ctx context.Context, loanID int) (bool, error) {
	return false, nil
}

func setupLoanHandler(service loan_service.LoanService) (*LoanHandler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	h := NewLoanHandler(service)
	r := gin.New()

	r.POST("/api/v1/loans", h.CreateLoan)

	return h, r
}

func TestLoanHandler_CreateLoan_Success(t *testing.T) {
	m := &mockLoanService{
		createResult: &model.Loan{
			ID:         1,
			BorrowerID: 1,
		},
	}
	_, r := setupLoanHandler(m)

	body := map[string]float64{
		"borrower_id": 1,
		"amount":      5000000,
	}
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal body: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/loans", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var resp model.Loan
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.ID != 1 || resp.BorrowerID != 1 {
		t.Fatalf("unexpected loan response: %+v", resp)
	}
}

func TestLoanHandler_CreateLoan_InvalidBody(t *testing.T) {
	m := &mockLoanService{}
	_, r := setupLoanHandler(m)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/loans", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
