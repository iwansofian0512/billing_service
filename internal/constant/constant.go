package constant

import "time"

const (
	DefaultPort     = "8080"
	ShutdownTimeout = 30 * time.Second
	MaxLoanDuration = 50
	LoanInterest    = 0.10
)
