package admin

import (
	"BookMyTurf/db"
	"BookMyTurf/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Blocking users

func BlockUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := db.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Blocked = true
	if err := db.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to block user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User blocked successfully"})
}

// Unblocking users

func UnblockUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := db.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Blocked = false
	if err := db.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unblock user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User unblocked successfully"})
}

// Get all users

func GetAllUsers(c *gin.Context) {
	var users []models.User
	if err := db.DB.Where("role = ?", 0).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func AdminUpdateUser(c *gin.Context) {
	// Bind JSON input with user_id
	var input struct {
		UserID uint   `json:"user_id" binding:"required"`
		Name   string `json:"name" binding:"required,min=2"`
		Email  string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	// Fetch user to be updated
	var user models.User
	if err := db.DB.First(&user, input.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update fields
	user.Name = input.Name
	user.Email = input.Email

	if err := db.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update user",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}
