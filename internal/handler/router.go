package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/iwansofian0512/billing_service/internal/handler/borrower_handler"
	"github.com/iwansofian0512/billing_service/internal/handler/loan_handler"
	"github.com/iwansofian0512/billing_service/internal/handler/payment_handler"
)

func NewRouter(loanHandler *loan_handler.LoanHandler, borrowerHandler *borrower_handler.BorrowerHandler, paymentHandler *payment_handler.PaymentHandler) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")

	// BORROWER
	api.POST("/borrowers", borrowerHandler.CreateBorrower)
	api.GET("/borrowers", borrowerHandler.ListBorrowerLoans)

	// LOAN
	api.POST("/loans", loanHandler.CreateLoan)
	api.POST("/payment", paymentHandler.MakePayment)

	return r
}
