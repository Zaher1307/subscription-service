package handlers

import (
	"net/http"
	"strconv"

	"github.com/zaher1307/subscription-service/internal/services"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	subscriptionService services.ISubscriptionService
}

func NewSubscriptionHandler(subscriptionService services.ISubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{subscriptionService: subscriptionService}
}

func (h *SubscriptionHandler) Create(c *gin.Context) {
	var request struct {
		UserID    int `json:"user_id" binding:"required"`
		ProductID int `json:"product_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscription, bill, err := h.subscriptionService.CreateSubscription(request.UserID, request.ProductID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"subscription": subscription,
		"initial_bill": bill,
	})
}

func (h *SubscriptionHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	subscription, err := h.subscriptionService.GetSubscription(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	c.JSON(http.StatusOK, subscription)
}
