package handlers

import (
	"net/http"
	"strconv"

	"github.com/zaher1307/subscription-service/internal/services"

	"github.com/gin-gonic/gin"
)

type BillHandler struct {
	billingService services.IBillingService
}

func NewBillHandler(billingService services.IBillingService) *BillHandler {
	return &BillHandler{billingService: billingService}
}

func (h *BillHandler) GetUserBills(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	bills, err := h.billingService.GetUserBills(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bills)
}

func (h *BillHandler) GetBill(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	bill, err := h.billingService.GetBill(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bill not found"})
		return
	}

	c.JSON(http.StatusOK, bill)
}

func (h *BillHandler) PayBill(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.billingService.PayBill(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bill paid successfully"})
}
