package auth

import (
	"BookMyTurf/config"
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

	err := config.DB.Where("email = ?", input.Email).First(&user).Error

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

	// Generate Refresh Token

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Name, user.Role)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "access token generation failed",
		})
		return
	}

	// Generating Refresh Token

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	config.DB.Create(&models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "refresh token generation failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Message":       "Login successful",
		"Access Token":  accessToken,
		"Refresh Token": refreshToken,
		"Role":          user.Role,
	})
}
