package services

import (
	"github.com/zaher1307/subscription-service/internal/models"
	"github.com/zaher1307/subscription-service/internal/repositories"
)

type ProductService struct {
	productRepo repositories.IProductRepository
}

func NewProductService(productRepo repositories.IProductRepository) *ProductService {
	return &ProductService{productRepo: productRepo}
}

func (s *ProductService) GetAllProducts() ([]*models.Product, error) {
	return s.productRepo.GetAll()
}

func (s *ProductService) GetProductByID(id int) (*models.Product, error) {
	return s.productRepo.GetByID(id)
}
