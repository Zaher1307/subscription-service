package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/zaher1307/subscription-service/internal/models"
)

type BillRepository struct {
	DB *sql.DB
}

func NewBillRepository(db *sql.DB) *BillRepository {
	return &BillRepository{DB: db}
}

func (r *BillRepository) Create(bill *models.Bill) error {
	stmt, err := r.DB.Prepare(`
		INSERT INTO bills (subscription_id, amount, status, paid_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.QueryRow(
		bill.SubscriptionID,
		bill.Amount,
		bill.Status,
		bill.PaidAt,
	).Scan(&bill.ID, &bill.CreatedAt)
}

func (r *BillRepository) GetByID(id int) (*models.Bill, error) {
	stmt, err := r.DB.Prepare(`
		SELECT id, subscription_id, amount, status, created_at, paid_at
		FROM bills
		WHERE id = $1
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var bill models.Bill
	err = stmt.QueryRow(id).Scan(
		&bill.ID,
		&bill.SubscriptionID,
		&bill.Amount,
		&bill.Status,
		&bill.CreatedAt,
		&bill.PaidAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("bill %d not found", id)
		}
		return nil, err
	}

	return &bill, nil
}

func (r *BillRepository) GetByUserID(userID int) ([]*models.Bill, error) {
	stmt, err := r.DB.Prepare(`
		SELECT b.id, b.subscription_id, b.amount, b.status, b.created_at, b.paid_at
		FROM bills b
		JOIN subscriptions s ON b.subscription_id = s.id
		WHERE s.user_id = $1
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bills := make([]*models.Bill, 0)
	for rows.Next() {
		var bill models.Bill
		if err := rows.Scan(
			&bill.ID,
			&bill.SubscriptionID,
			&bill.Amount,
			&bill.Status,
			&bill.CreatedAt,
			&bill.PaidAt,
		); err != nil {
			return nil, err
		}
		bills = append(bills, &bill)
	}

	return bills, nil
}

func (r *BillRepository) MarkAsPaid(id int) error {
	stmt, err := r.DB.Prepare(`
		UPDATE bills
		SET status = 'paid', paid_at = $1
		WHERE id = $2
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	now := time.Now()
	result, err := stmt.Exec(now, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
