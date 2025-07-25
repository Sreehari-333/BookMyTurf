package payment

import (
	"BookMyTurf/db"
	"BookMyTurf/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreatePayment(c *gin.Context) {
	var payment models.Payment

	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	if payment.UserID == 0 || payment.BookingID == 0 || payment.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid fields"})
		return
	}

	payment.Status = "Success"
	if err := db.DB.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Payment recorded", "payment": payment})
}

func GetAllPayments(c *gin.Context) {
	var payments []models.Payment
	if err := db.DB.Find(&payments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payments"})
		return
	}
	c.JSON(http.StatusOK, payments)
}

func GetPaymentByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var payment models.Payment
	if err := db.DB.First(&payment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	c.JSON(http.StatusOK, payment)
}

func DeletePayment(c *gin.Context) {
	id := c.Param("id")
	var payment models.Payment
	if err := db.DB.First(&payment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}
	if err := db.DB.Delete(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete payment"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Payment deleted"})
}
