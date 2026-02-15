package loan_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iwansofian0512/billing_service/internal/model"
	"github.com/iwansofian0512/billing_service/internal/service/loan_service"
)

type LoanHandler struct {
	service loan_service.LoanService
}

func NewLoanHandler(service loan_service.LoanService) *LoanHandler {
	return &LoanHandler{service: service}
}

func (h *LoanHandler) CreateLoan(ctx *gin.Context) {
	var req model.CreateLoanRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loan, err := h.service.CreateLoan(ctx, int(req.BorrowerID), req.Amount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	loan.Schedules = nil

	ctx.JSON(http.StatusCreated, loan)
}
