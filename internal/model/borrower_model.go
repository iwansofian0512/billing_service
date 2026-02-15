package model

import "time"

type CreateBorrowerRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Borrower struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email,omitempty" db:"email"`
	IsActive  bool      `json:"isActive" db:"is_active"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}
