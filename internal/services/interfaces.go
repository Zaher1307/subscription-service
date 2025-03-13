package services

import "github.com/zaher1307/subscription-service/internal/models"

type ISubscriptionService interface {
	CreateSubscription(userID, productID int) (*models.Subscription, *models.Bill, error)
	GetSubscription(id int) (*models.Subscription, error)
}

var _ ISubscriptionService = (*SubscriptionService)(nil)

type IBillingService interface {
	GetBill(id int) (*models.Bill, error)
	GetUserBills(userID int) ([]*models.Bill, error)
	PayBill(id int) error
	GenerateBills() error
}

var _ IBillingService = (*BillingService)(nil)

type IProductService interface {
	GetAllProducts() ([]*models.Product, error)
	GetProductByID(id int) (*models.Product, error)
}

var _ IProductService = (*ProductService)(nil)

type IUserService interface {
	CreateUser(user *models.User) error
	GetUserByID(id int) (*models.User, error)
}

var _ IUserService = (*UserService)(nil)
