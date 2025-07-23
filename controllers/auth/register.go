package auth

import (
	"BookMyTurf/db"
	"BookMyTurf/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Register Handler

func Register(c *gin.Context) {

	var input models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Hashing password

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to hash password",
		})
		return
	}

	// role := 0

	// if input.Role != nil {
	// 	role = *input.Role
	// }

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     input.Role,
	}

	// creating user

	err = db.DB.Create(&user).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "user creation error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Message": "Registration successful",
	})

}
