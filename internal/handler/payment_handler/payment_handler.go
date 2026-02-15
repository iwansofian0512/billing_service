package payment_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iwansofian0512/billing_service/internal/model"
	"github.com/iwansofian0512/billing_service/internal/service/payment_service"
)

type PaymentHandler struct {
	service payment_service.PaymentService
}

func NewPaymentHandler(service payment_service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

func (h *PaymentHandler) MakePayment(ctx *gin.Context) {
	var req model.PaymentRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.LoanID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid loanID"})
		return
	}

	err := h.service.MakePayment(ctx, req.LoanID, req.Amount)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "payment successful"})
}
