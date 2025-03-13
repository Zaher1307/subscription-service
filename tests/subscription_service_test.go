package tests

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zaher1307/subscription-service/internal/models"
	"github.com/zaher1307/subscription-service/internal/repositories"
	"github.com/zaher1307/subscription-service/internal/services"
)

type MockSubscriptionRepository struct {
	mock.Mock
}

var _ repositories.ISubscriptionRepository = (*MockSubscriptionRepository)(nil)

func (m *MockSubscriptionRepository) Create(subscription *models.Subscription) error {
	args := m.Called(subscription)
	subscription.ID = 1
	subscription.CreatedAt = time.Now()
	return args.Error(0)
}

func (m *MockSubscriptionRepository) GetByID(id int) (*models.Subscription, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) GetActiveByUserAndProduct(userID, productID int) (*models.Subscription, error) {
	args := m.Called(userID, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) GetDueForBilling(date time.Time) ([]*models.Subscription, error) {
	args := m.Called(date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) UpdateNextBillingDate(id int, nextDate time.Time) error {
	args := m.Called(id, nextDate)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) UpdateStartDate(id int, nextDate time.Time) error {
	args := m.Called(id, nextDate)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) HoldSubscription(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) ActivateSubscription(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockProductRepository struct {
	mock.Mock
}

var _ repositories.IProductRepository = (*MockProductRepository)(nil)

func (m *MockProductRepository) GetByID(id int) (*models.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepository) GetAll() ([]*models.Product, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Product), args.Error(1)
}

type MockBillRepository struct {
	mock.Mock
}

var _ repositories.IBillRepository = (*MockBillRepository)(nil)

func (m *MockBillRepository) Create(bill *models.Bill) error {
	args := m.Called(bill)
	bill.ID = 1
	return args.Error(0)
}

func (m *MockBillRepository) GetByID(id int) (*models.Bill, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Bill), args.Error(1)
}

func (m *MockBillRepository) GetByUserID(userID int) ([]*models.Bill, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Bill), args.Error(1)
}

func (m *MockBillRepository) MarkAsPaid(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockUserRepository struct {
	mock.Mock
}

var _ repositories.IUserRepository = (*MockUserRepository)(nil)

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	user.ID = 1
	user.CreatedAt = time.Now()
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id int) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestSubscriptionService_CreateSubscription(t *testing.T) {
	tests := []struct {
		name                  string
		userID                int
		productID             int
		mockSetup             func(mockSubRepo *MockSubscriptionRepository, mockProductRepo *MockProductRepository, mockBillRepo *MockBillRepository, mockUserRepo *MockUserRepository)
		expectedSubscription  *models.Subscription
		expectedBill          *models.Bill
		expectedError         bool
		expectedErrorContains string
	}{
		{
			name:      "successful subscription creation",
			userID:    1,
			productID: 2,
			mockSetup: func(mockSubRepo *MockSubscriptionRepository, mockProductRepo *MockProductRepository, mockBillRepo *MockBillRepository, mockUserRepo *MockUserRepository) {
				mockSubRepo.On("GetActiveByUserAndProduct", 1, 2).Return(nil, nil)

				mockUserRepo.On("GetByID", 1).Return(&models.User{ID: 1, Name: "Test User"}, nil)

				mockProductRepo.On("GetByID", 2).Return(&models.Product{ID: 2, Name: "Test Product", Price: 99.99}, nil)

				mockSubRepo.On("Create", mock.AnythingOfType("*models.Subscription")).Return(nil)

				mockBillRepo.On("Create", mock.AnythingOfType("*models.Bill")).Return(nil)
			},
			expectedSubscription: &models.Subscription{
				ID:        1,
				UserID:    1,
				ProductID: 2,
				Status:    "active",
			},
			expectedBill: &models.Bill{
				ID:             1,
				SubscriptionID: 1,
				Amount:         99.99,
				Status:         "paid",
			},
			expectedError: false,
		},
		{
			name:      "user already has subscription",
			userID:    1,
			productID: 2,
			mockSetup: func(mockSubRepo *MockSubscriptionRepository, mockProductRepo *MockProductRepository, mockBillRepo *MockBillRepository, mockUserRepo *MockUserRepository) {
				existingSub := &models.Subscription{
					ID:        1,
					UserID:    1,
					ProductID: 2,
					Status:    "active",
				}
				mockSubRepo.On("GetActiveByUserAndProduct", 1, 2).Return(existingSub, nil)
			},
			expectedSubscription:  nil,
			expectedBill:          nil,
			expectedError:         true,
			expectedErrorContains: "user already has an active subscription",
		},
		{
			name:      "user not found",
			userID:    999,
			productID: 2,
			mockSetup: func(mockSubRepo *MockSubscriptionRepository, mockProductRepo *MockProductRepository, mockBillRepo *MockBillRepository, mockUserRepo *MockUserRepository) {
				mockSubRepo.On("GetActiveByUserAndProduct", 999, 2).Return(nil, nil)
				mockUserRepo.On("GetByID", 999).Return(nil, errors.New("user not found"))
			},
			expectedSubscription:  nil,
			expectedBill:          nil,
			expectedError:         true,
			expectedErrorContains: "user not found",
		},
		{
			name:      "product not found",
			userID:    1,
			productID: 999,
			mockSetup: func(mockSubRepo *MockSubscriptionRepository, mockProductRepo *MockProductRepository, mockBillRepo *MockBillRepository, mockUserRepo *MockUserRepository) {
				mockSubRepo.On("GetActiveByUserAndProduct", 1, 999).Return(nil, nil)
				mockUserRepo.On("GetByID", 1).Return(&models.User{ID: 1, Name: "Test User"}, nil)
				mockProductRepo.On("GetByID", 999).Return(nil, errors.New("product not found"))
			},
			expectedSubscription:  nil,
			expectedBill:          nil,
			expectedError:         true,
			expectedErrorContains: "product not found",
		},
		{
			name:      "subscription creation fails",
			userID:    1,
			productID: 2,
			mockSetup: func(mockSubRepo *MockSubscriptionRepository, mockProductRepo *MockProductRepository, mockBillRepo *MockBillRepository, mockUserRepo *MockUserRepository) {
				mockSubRepo.On("GetActiveByUserAndProduct", 1, 2).Return(nil, nil)
				mockUserRepo.On("GetByID", 1).Return(&models.User{ID: 1, Name: "Test User"}, nil)
				mockProductRepo.On("GetByID", 2).Return(&models.Product{ID: 2, Name: "Test Product", Price: 99.99}, nil)
				mockSubRepo.On("Create", mock.AnythingOfType("*models.Subscription")).Return(errors.New("db error"))
			},
			expectedSubscription:  nil,
			expectedBill:          nil,
			expectedError:         true,
			expectedErrorContains: "db error",
		},
		{
			name:      "bill creation fails",
			userID:    1,
			productID: 2,
			mockSetup: func(mockSubRepo *MockSubscriptionRepository, mockProductRepo *MockProductRepository, mockBillRepo *MockBillRepository, mockUserRepo *MockUserRepository) {
				mockSubRepo.On("GetActiveByUserAndProduct", 1, 2).Return(nil, nil)
				mockUserRepo.On("GetByID", 1).Return(&models.User{ID: 1, Name: "Test User"}, nil)
				mockProductRepo.On("GetByID", 2).Return(&models.Product{ID: 2, Name: "Test Product", Price: 99.99}, nil)
				mockSubRepo.On("Create", mock.AnythingOfType("*models.Subscription")).Return(nil)
				mockBillRepo.On("Create", mock.AnythingOfType("*models.Bill")).Return(errors.New("db error"))
			},
			expectedSubscription:  nil,
			expectedBill:          nil,
			expectedError:         true,
			expectedErrorContains: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSubscriptionRepo := new(MockSubscriptionRepository)
			mockProductRepo := new(MockProductRepository)
			mockBillRepo := new(MockBillRepository)
			mockUserRepo := new(MockUserRepository)

			tt.mockSetup(mockSubscriptionRepo, mockProductRepo, mockBillRepo, mockUserRepo)

			service := services.NewSubscriptionService(
				mockSubscriptionRepo,
				mockProductRepo,
				mockBillRepo,
				mockUserRepo,
			)

			subscription, bill, err := service.CreateSubscription(tt.userID, tt.productID)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrorContains != "" {
					assert.Contains(t, err.Error(), tt.expectedErrorContains)
				}
				assert.Nil(t, subscription)
				assert.Nil(t, bill)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, subscription)
				assert.NotNil(t, bill)
				assert.Equal(t, tt.expectedSubscription.UserID, subscription.UserID)
				assert.Equal(t, tt.expectedSubscription.ProductID, subscription.ProductID)
				assert.Equal(t, tt.expectedSubscription.Status, subscription.Status)
				assert.Equal(t, tt.expectedBill.Amount, bill.Amount)
				assert.Equal(t, tt.expectedBill.Status, bill.Status)
			}

			mockSubscriptionRepo.AssertExpectations(t)
			mockProductRepo.AssertExpectations(t)
			mockBillRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestSubscriptionService_GetSubscription(t *testing.T) {
	tests := []struct {
		name                  string
		subscriptionID        int
		mockSetup             func(mockRepo *MockSubscriptionRepository)
		expectedSubscription  *models.Subscription
		expectedError         bool
		expectedErrorContains string
	}{
		{
			name:           "subscription found",
			subscriptionID: 1,
			mockSetup: func(mockRepo *MockSubscriptionRepository) {
				now := time.Now()
				nextMonth := now.AddDate(0, 1, 0)
				subscription := &models.Subscription{
					ID:              1,
					UserID:          1,
					ProductID:       2,
					StartDate:       now,
					NextBillingDate: nextMonth,
					Status:          "active",
					CreatedAt:       now,
				}
				mockRepo.On("GetByID", 1).Return(subscription, nil)
			},
			expectedSubscription: &models.Subscription{
				ID:     1,
				UserID: 1,
				Status: "active",
			},
			expectedError: false,
		},
		{
			name:           "subscription not found",
			subscriptionID: 999,
			mockSetup: func(mockRepo *MockSubscriptionRepository) {
				mockRepo.On("GetByID", 999).Return(nil, errors.New("Subscription 999 not found"))
			},
			expectedSubscription:  nil,
			expectedError:         true,
			expectedErrorContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockSubscriptionRepository)

			tt.mockSetup(mockRepo)

			mockProductRepo := new(MockProductRepository)
			mockBillRepo := new(MockBillRepository)
			mockUserRepo := new(MockUserRepository)

			service := services.NewSubscriptionService(mockRepo, mockProductRepo, mockBillRepo, mockUserRepo)

			subscription, err := service.GetSubscription(tt.subscriptionID)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrorContains != "" {
					assert.Contains(t, err.Error(), tt.expectedErrorContains)
				}
				assert.Nil(t, subscription)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, subscription)
				assert.Equal(t, tt.expectedSubscription.ID, subscription.ID)
				assert.Equal(t, tt.expectedSubscription.UserID, subscription.UserID)
				assert.Equal(t, tt.expectedSubscription.Status, subscription.Status)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
