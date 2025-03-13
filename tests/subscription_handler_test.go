package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zaher1307/subscription-service/internal/handlers"
	"github.com/zaher1307/subscription-service/internal/models"
	"github.com/zaher1307/subscription-service/internal/services"
)

type MockSubscriptionService struct {
	mock.Mock
}

var _ services.ISubscriptionService = (*MockSubscriptionService)(nil)

func (m *MockSubscriptionService) CreateSubscription(userID, productID int) (*models.Subscription, *models.Bill, error) {
	args := m.Called(userID, productID)
	subscription, _ := args.Get(0).(*models.Subscription)
	bill, _ := args.Get(1).(*models.Bill)
	return subscription, bill, args.Error(2)
}

func (m *MockSubscriptionService) GetSubscription(id int) (*models.Subscription, error) {
	args := m.Called(id)
	subscription, _ := args.Get(0).(*models.Subscription)
	return subscription, args.Error(1)
}

func TestSubscriptionHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		requestBody        map[string]interface{}
		mockSetup          func(mock *MockSubscriptionService)
		expectedStatusCode int
		expectedResponse   map[string]interface{}
	}{
		{
			name: "valid request",
			requestBody: map[string]interface{}{
				"user_id":    1,
				"product_id": 2,
			},
			mockSetup: func(mockService *MockSubscriptionService) {
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

				bill := &models.Bill{
					ID:             1,
					SubscriptionID: 1,
					Amount:         99.99,
					Status:         "paid",
					PaidAt:         &now,
				}

				mockService.On("CreateSubscription", 1, 2).Return(subscription, bill, nil)
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse: map[string]interface{}{
				"subscription": mock.Anything,
				"initial_bill": mock.Anything,
			},
		},
		{
			name: "missing required fields",
			requestBody: map[string]interface{}{
				"user_id": 1,
				// missing product_id
			},
			mockSetup:          func(mockService *MockSubscriptionService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": mock.Anything,
			},
		},
		{
			name: "service error",
			requestBody: map[string]interface{}{
				"user_id":    1,
				"product_id": 2,
			},
			mockSetup: func(mockService *MockSubscriptionService) {
				mockService.On("CreateSubscription", 1, 2).Return(nil, nil, errors.New("service error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{
				"error": "service error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSubscriptionService)
			tt.mockSetup(mockService)

			handler := handlers.NewSubscriptionHandler(mockService)

			router := gin.New()
			router.POST("/subscriptions", handler.Create)

			requestJSON, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/subscriptions", bytes.NewBuffer(requestJSON))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			var response map[string]interface{}
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)

			for key, expectedValue := range tt.expectedResponse {
				assert.Contains(t, response, key)
				if expectedValue != mock.Anything {
					assert.Equal(t, expectedValue, response[key])
				}
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestSubscriptionHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		subscriptionID     string
		mockSetup          func(mock *MockSubscriptionService)
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name:           "valid request",
			subscriptionID: "1",
			mockSetup: func(mockService *MockSubscriptionService) {
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

				mockService.On("GetSubscription", 1).Return(subscription, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   mock.Anything,
		},
		{
			name:               "invalid id",
			subscriptionID:     "invalid",
			mockSetup:          func(mockService *MockSubscriptionService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "Invalid ID",
			},
		},
		{
			name:           "not found",
			subscriptionID: "999",
			mockSetup: func(mockService *MockSubscriptionService) {
				mockService.On("GetSubscription", 999).Return(nil, errors.New("Subscription not found"))
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: map[string]interface{}{
				"error": "Subscription not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSubscriptionService)
			tt.mockSetup(mockService)

			handler := handlers.NewSubscriptionHandler(mockService)

			router := gin.New()
			router.GET("/subscriptions/:id", handler.GetByID)

			req, _ := http.NewRequest(http.MethodGet, "/subscriptions/"+tt.subscriptionID, nil)

			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			if mapResponse, ok := tt.expectedResponse.(map[string]interface{}); ok {
				var response map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)

				for key, expectedValue := range mapResponse {
					assert.Contains(t, response, key)
					if expectedValue != mock.Anything {
						assert.Equal(t, expectedValue, response[key])
					}
				}
			}

			mockService.AssertExpectations(t)
		})
	}
}
