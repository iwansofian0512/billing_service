package borrower_handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/iwansofian0512/billing_service/internal/model"
	"github.com/iwansofian0512/billing_service/internal/service/borrower_service"
)

type mockBorrowerService struct {
	createResult *model.Borrower
	createErr    error
	listResult   []model.Loan
	listErr      error
}

func (m *mockBorrowerService) CreateBorrower(ctx context.Context, name, email string) (*model.Borrower, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	if m.createResult != nil {
		return m.createResult, nil
	}
	return &model.Borrower{
		ID:    1,
		Name:  name,
		Email: email,
	}, nil
}

func (m *mockBorrowerService) GetBorrowerLoans(ctx context.Context) ([]model.Loan, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.listResult, nil
}

func (m *mockBorrowerService) ListBorrowerLoans(ctx context.Context, borrowerID, page, pageSize int) ([]model.Loan, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.listResult, nil
}

func setupBorrowerHandler(service borrower_service.BorrowerService) (*BorrowerHandler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	h := NewBorrowerHandler(service)
	r := gin.New()
	r.POST("/api/v1/borrowers", h.CreateBorrower)
	r.GET("/api/v1/borrowers", h.ListBorrowerLoans)
	return h, r
}

func TestBorrowerHandler_CreateBorrower_Success(t *testing.T) {
	m := &mockBorrowerService{}
	_, r := setupBorrowerHandler(m)

	body := map[string]string{
		"name":  "John Doe",
		"email": "john@example.com",
	}
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal body: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/borrowers", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var resp model.Borrower
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Name != "John Doe" || resp.Email != "john@example.com" {
		t.Fatalf("unexpected borrower response: %+v", resp)
	}
}

func TestBorrowerHandler_CreateBorrower_InvalidBody(t *testing.T) {
	m := &mockBorrowerService{}
	_, r := setupBorrowerHandler(m)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/borrowers", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
