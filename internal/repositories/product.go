package repositories

import (
	"database/sql"
	"fmt"

	"github.com/zaher1307/subscription-service/internal/models"
)

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) GetAll() ([]*models.Product, error) {
	stmt, err := r.DB.Prepare(`
		SELECT id, name, description, price, created_at
		FROM products
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]*models.Product, 0)
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.CreatedAt); err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, nil
}

func (r *ProductRepository) GetByID(id int) (*models.Product, error) {
	stmt, err := r.DB.Prepare(`
		SELECT id, name, description, price, created_at
		FROM products
		WHERE id = $1
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var product models.Product
	err = stmt.QueryRow(id).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product %d not found", id)
		}
		return nil, err
	}

	return &product, nil
}
