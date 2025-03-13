package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/zaher1307/subscription-service/internal/models"
)

type SubscriptionRepository struct {
	DB *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{DB: db}
}

func (r *SubscriptionRepository) Create(subscription *models.Subscription) error {
	stmt, err := r.DB.Prepare(`
		INSERT INTO subscriptions (user_id, product_id, start_date, next_billing_date, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.QueryRow(
		subscription.UserID,
		subscription.ProductID,
		subscription.StartDate,
		subscription.NextBillingDate,
		subscription.Status,
	).Scan(&subscription.ID, &subscription.CreatedAt)
}

func (r *SubscriptionRepository) GetByID(id int) (*models.Subscription, error) {
	stmt, err := r.DB.Prepare(`
		SELECT id, user_id, product_id, start_date, next_billing_date, status, created_at
		FROM subscriptions
		WHERE id = $1
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var subscription models.Subscription
	err = stmt.QueryRow(id).Scan(
		&subscription.ID,
		&subscription.UserID,
		&subscription.ProductID,
		&subscription.StartDate,
		&subscription.NextBillingDate,
		&subscription.Status,
		&subscription.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("subscription %d not found", id)
		}
		return nil, err
	}

	return &subscription, nil
}

func (r *SubscriptionRepository) GetDueForBilling(date time.Time) ([]*models.Subscription, error) {
	stmt, err := r.DB.Prepare(`
		SELECT id, user_id, product_id, start_date, next_billing_date, status, created_at
		FROM subscriptions
		WHERE status = 'active' AND next_billing_date <= $1
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subscriptions := make([]*models.Subscription, 0)
	for rows.Next() {
		var subscription models.Subscription
		if err := rows.Scan(
			&subscription.ID,
			&subscription.UserID,
			&subscription.ProductID,
			&subscription.StartDate,
			&subscription.NextBillingDate,
			&subscription.Status,
			&subscription.CreatedAt,
		); err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, &subscription)
	}

	return subscriptions, nil
}

func (r *SubscriptionRepository) UpdateStartDate(id int, nextDate time.Time) error {
	stmt, err := r.DB.Prepare(`
		UPDATE subscriptions
		SET start_date = $1
		WHERE id = $2
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(nextDate, id)
	return err
}

func (r *SubscriptionRepository) UpdateNextBillingDate(id int, nextDate time.Time) error {
	stmt, err := r.DB.Prepare(`
		UPDATE subscriptions
		SET next_billing_date = $1
		WHERE id = $2
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(nextDate, id)
	return err
}

func (r *SubscriptionRepository) HoldSubscription(id int) error {
	stmt, err := r.DB.Prepare(`
		UPDATE subscriptions
		SET status = 'hold'
		WHERE id = $1
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

func (r *SubscriptionRepository) ActivateSubscription(id int) error {
	stmt, err := r.DB.Prepare(`
		UPDATE subscriptions
		SET status = 'active'
		WHERE id = $1
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

func (r *SubscriptionRepository) GetActiveByUserAndProduct(userID, productID int) (*models.Subscription, error) {
	stmt, err := r.DB.Prepare(`
        SELECT id, user_id, product_id, start_date, next_billing_date, status, created_at
        FROM subscriptions
        WHERE user_id = $1 AND product_id = $2 AND status = 'active'
    `)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var subscription models.Subscription
	err = stmt.QueryRow(userID, productID).Scan(
		&subscription.ID,
		&subscription.UserID,
		&subscription.ProductID,
		&subscription.StartDate,
		&subscription.NextBillingDate,
		&subscription.Status,
		&subscription.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &subscription, nil
}
