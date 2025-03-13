package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/zaher1307/subscription-service/internal/models"
	"github.com/zaher1307/subscription-service/internal/repositories"
)

type SubscriptionService struct {
	subscriptionRepo repositories.ISubscriptionRepository
	productRepo      repositories.IProductRepository
	billRepo         repositories.IBillRepository
	userRepo         repositories.IUserRepository
}

func NewSubscriptionService(
	subscriptionRepo repositories.ISubscriptionRepository,
	productRepo repositories.IProductRepository,
	billRepo repositories.IBillRepository,
	userRepo repositories.IUserRepository,
) *SubscriptionService {
	return &SubscriptionService{
		subscriptionRepo: subscriptionRepo,
		productRepo:      productRepo,
		billRepo:         billRepo,
		userRepo:         userRepo,
	}
}

func (s *SubscriptionService) CreateSubscription(userID, productID int) (*models.Subscription, *models.Bill, error) {
	existingSub, err := s.subscriptionRepo.GetActiveByUserAndProduct(userID, productID)
	if err != nil {
		return nil, nil, err
	}
	if existingSub != nil {
		return nil, nil, errors.New("user already has an active subscription for this product")
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, nil, err
	}

	product, err := s.productRepo.GetByID(productID)
	if err != nil {
		return nil, nil, err
	}

	now := time.Now()
	nextMonth := now.AddDate(0, 1, 0)

	subscription := &models.Subscription{
		UserID:          user.ID,
		ProductID:       product.ID,
		StartDate:       now,
		NextBillingDate: nextMonth,
		Status:          "active",
	}

	if err := s.subscriptionRepo.Create(subscription); err != nil {
		return nil, nil, err
	}

	bill := &models.Bill{
		SubscriptionID: subscription.ID,
		Amount:         product.Price,
		Status:         "paid",
		PaidAt:         &now,
	}

	fmt.Printf("bill: %#v\n", bill)

	if err := s.billRepo.Create(bill); err != nil {
		return nil, nil, err
	}

	return subscription, bill, nil
}

func (s *SubscriptionService) GetSubscription(id int) (*models.Subscription, error) {
	return s.subscriptionRepo.GetByID(id)
}
