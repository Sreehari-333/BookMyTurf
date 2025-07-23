package auth

import (
	"BookMyTurf/db"
	"BookMyTurf/models"
	"BookMyTurf/utils"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	jwtLib "github.com/golang-jwt/jwt/v5"
)

var secret = os.Getenv("JWT_SECRET")

func RefreshToken(c *gin.Context) {

	//  Reading token from cookie

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token missing"})
		return
	}

	//  validating token

	token, err := jwtLib.ParseWithClaims(refreshToken, &utils.Claims{}, func(token *jwtLib.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	claims, ok := token.Claims.(*utils.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	// Checking if token exists in DB

	var existingToken models.RefreshToken
	err = db.DB.Where("user_id = ? AND token = ?", claims.UserId, refreshToken).First(&existingToken).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not recognized"})
		return
	}

	// Deleting existing token from DB

	db.DB.Delete(&existingToken)

	// Generating new Access Token

	newAccessToken, err := utils.GenerateAccessToken(claims.UserId, claims.UserName, claims.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create access token"})
		return
	}

	// Generating new Refresh Token

	newRefreshToken, err := utils.GenerateRefreshToken(claims.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create refresh token"})
		return
	}

	//  Store new refresh token in DB

	db.DB.Create(&models.RefreshToken{
		UserID:    claims.UserId,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	})

	c.SetCookie("refresh_token", newRefreshToken, 7*24*3600, "/", "localhost", false, true) // Setting new cookie

	// Return new access token

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
	})
}
