package utils

import (
	jwtLib "github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type Claims struct {
	UserId   uint
	UserName string
	Role     int
	jwtLib.RegisteredClaims
}

var secret = []byte(os.Getenv("JWT_SECRET")) // secret key

// Generating Access Token

func GenerateAccessToken(userId uint, userName string, role int) (string, error) {

	expirationTime := time.Now().Add(30 * time.Minute)

	claims := Claims{
		UserId:   userId,
		UserName: userName,
		Role:     role,
		RegisteredClaims: jwtLib.RegisteredClaims{
			ExpiresAt: jwtLib.NewNumericDate(expirationTime),
			IssuedAt:  jwtLib.NewNumericDate(time.Now()),
		},
	}

	token := jwtLib.NewWithClaims(jwtLib.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

// Generating Refresh Token

func GenerateRefreshToken(userId uint) (string, error) {

	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	claims := Claims{
		UserId: userId,
		RegisteredClaims: jwtLib.RegisteredClaims{
			ExpiresAt: jwtLib.NewNumericDate(expirationTime),
			IssuedAt:  jwtLib.NewNumericDate(time.Now()),
		},
	}

	token := jwtLib.NewWithClaims(jwtLib.SigningMethodHS256, claims)

	return token.SignedString(secret)
}
