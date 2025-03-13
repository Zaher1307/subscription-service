package models

import "time"

type Bill struct {
	ID             int        `json:"id"`
	SubscriptionID int        `json:"subscription_id"`
	Amount         float64    `json:"amount"`
	Status         string     `json:"status"`
	CreatedAt      time.Time  `json:"created_at"`
	PaidAt         *time.Time `json:"paid_at,omitempty"`
}
