package auth

import (
	"BookMyTurf/db"
	"BookMyTurf/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {

	//  Reading Refresh Token from cookie

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token missing in cookie"})
		return
	}

	//  Delete Refresh Token from DB

	err = db.DB.Where("token = ?", refreshToken).Delete(&models.RefreshToken{}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true) //  Clearing cookie

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
