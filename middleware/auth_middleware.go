package middleware

import (
	"BookMyTurf/constants"
	"BookMyTurf/utils"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// struct for extracting data from token

var secret = []byte(os.Getenv("JWT_SECRET")) // secret key

// Auth Middleware

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing header",
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims := &utils.Claims{}

		// Parsing claims

		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return secret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or token expired",
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserId)
		c.Set("role", claims.Role)
		c.Next()

	}
}

// Role based access Middleware

func IsAdmin() gin.HandlerFunc {

	return func(c *gin.Context) {
		roleValue, exists := c.Get("role")

		fmt.Printf(" Role from context: %v (%T)\n", roleValue, roleValue)

		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "role not found",
			})
			c.Abort()
			return
		}

		role, ok := roleValue.(int)
		if !ok || role != int(constants.RoleAdmin) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "access denied admins only",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
