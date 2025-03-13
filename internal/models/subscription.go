package models

import "time"

type Subscription struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	ProductID       int       `json:"product_id"`
	StartDate       time.Time `json:"start_date"`
	NextBillingDate time.Time `json:"next_billing_date"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}
