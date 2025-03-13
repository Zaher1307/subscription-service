package services

import (
	"errors"
	"time"

	"github.com/zaher1307/subscription-service/internal/models"
	"github.com/zaher1307/subscription-service/internal/repositories"
)

type BillingService struct {
	subscriptionRepo repositories.ISubscriptionRepository
	productRepo      repositories.IProductRepository
	billRepo         repositories.IBillRepository
	userRepo         repositories.IUserRepository
}

func NewBillingService(
	subscriptionRepo repositories.ISubscriptionRepository,
	productRepo repositories.IProductRepository,
	billRepo repositories.IBillRepository,
	userRepo repositories.IUserRepository,
) *BillingService {
	return &BillingService{
		subscriptionRepo: subscriptionRepo,
		productRepo:      productRepo,
		billRepo:         billRepo,
		userRepo:         userRepo,
	}
}

func (s *BillingService) GetBill(id int) (*models.Bill, error) {
	return s.billRepo.GetByID(id)
}

func (s *BillingService) GetUserBills(userID int) ([]*models.Bill, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return s.billRepo.GetByUserID(user.ID)
}

func (s *BillingService) PayBill(id int) error {
	bill, err := s.billRepo.GetByID(id)
	if err != nil {
		return err
	}

	if bill.Status == "paid" {
		return errors.New("bill is already paid")
	}

	err = s.billRepo.MarkAsPaid(id)
	if err != nil {
		return err
	}

	now := time.Now()
	err = s.subscriptionRepo.UpdateStartDate(bill.SubscriptionID, now)
	if err != nil {
		return err
	}

	nextMonth := now.AddDate(0, 1, 0)
	err = s.subscriptionRepo.UpdateNextBillingDate(bill.SubscriptionID, nextMonth)
	if err != nil {
		return err
	}

	err = s.subscriptionRepo.ActivateSubscription(bill.SubscriptionID)
	if err != nil {
		return err
	}

	return nil
}

func (s *BillingService) GenerateBills() error {
	now := time.Now()

	subscriptions, err := s.subscriptionRepo.GetDueForBilling(now)
	if err != nil {
		return err
	}

	for _, subscription := range subscriptions {
		err := s.subscriptionRepo.HoldSubscription(subscription.ID)
		if err != nil {
			return err
		}

		product, err := s.productRepo.GetByID(subscription.ProductID)
		if err != nil {
			return err
		}

		bill := &models.Bill{
			SubscriptionID: subscription.ID,
			Amount:         product.Price,
			Status:         "pending",
		}

		if err := s.billRepo.Create(bill); err != nil {
			return err
		}
	}

	return nil
}
