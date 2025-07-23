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
