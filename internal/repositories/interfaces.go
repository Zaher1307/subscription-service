package repositories

import (
	"time"

	"github.com/zaher1307/subscription-service/internal/models"
)

type ISubscriptionRepository interface {
	Create(subscription *models.Subscription) error
	GetByID(id int) (*models.Subscription, error)
	GetActiveByUserAndProduct(userID, productID int) (*models.Subscription, error)
	GetDueForBilling(date time.Time) ([]*models.Subscription, error)
	UpdateNextBillingDate(id int, nextDate time.Time) error
	UpdateStartDate(id int, nextDate time.Time) error
	HoldSubscription(id int) error
	ActivateSubscription(id int) error
}

type IProductRepository interface {
	GetAll() ([]*models.Product, error)
	GetByID(id int) (*models.Product, error)
}

type IBillRepository interface {
	Create(bill *models.Bill) error
	GetByID(id int) (*models.Bill, error)
	GetByUserID(userID int) ([]*models.Bill, error)
	MarkAsPaid(id int) error
}

type IUserRepository interface {
	Create(user *models.User) error
	GetByID(id int) (*models.User, error)
}

var _ ISubscriptionRepository = (*SubscriptionRepository)(nil)
