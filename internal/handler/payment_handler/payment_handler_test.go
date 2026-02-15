package payment_handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/iwansofian0512/billing_service/internal/service/payment_service"
)

type mockPaymentService struct {
	err error
}

func (m *mockPaymentService) MakePayment(ctx context.Context, loanID int, amount float64) error {
	return m.err
}

func setupPaymentHandler(service payment_service.PaymentService) (*PaymentHandler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	h := NewPaymentHandler(service)
	r := gin.New()

	r.POST("/api/v1/payment", h.MakePayment)

	return h, r
}

func TestPaymentHandler_MakePayment_Success(t *testing.T) {
	m := &mockPaymentService{}
	_, r := setupPaymentHandler(m)

	body := map[string]interface{}{
		"loanID": 3,
		"amount": 110000,
	}
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal body: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/payment", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestPaymentHandler_MakePayment_InvalidID(t *testing.T) {
	m := &mockPaymentService{}
	_, r := setupPaymentHandler(m)

	body := map[string]interface{}{
		"loanID": 0,
		"amount": 110000,
	}
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal body: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/payment", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
