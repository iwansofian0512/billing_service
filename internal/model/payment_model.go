package model

type PaymentRequest struct {
	LoanID int     `json:"loanID"`
	Amount float64 `json:"amount"`
}
