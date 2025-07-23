package auth

import (
	"BookMyTurf/db"
	"BookMyTurf/models"
	"BookMyTurf/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Struct to store login credentials

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login

func Login(c *gin.Context) {

	var input LoginInput // login input

	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// checking the matching email

	var user models.User

	err := db.DB.Where("email = ?", input.Email).First(&user).Error

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// checking password

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid password",
		})
		return
	}

	// Invalidate old refresh tokens
	db.DB.Where("user_id = ?", user.ID).Delete(&models.RefreshToken{})

	// Generate tokens
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Name, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "access token generation failed"})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "refresh token generation failed"})
		return
	}

	// Store refresh token in DB
	db.DB.Create(&models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	})

	// Store refresh token securely in HttpOnly cookie
	c.SetCookie("refresh_token", refreshToken, 7*24*3600, "/", "localhost", false, true)

	// Send access token only
	c.JSON(http.StatusOK, gin.H{
		"message":       "Login successful",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"role":          user.Role,
	})

}
