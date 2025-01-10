package Config

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	expirationTime := time.Now().Add(1 * time.Hour).Unix()
	token.Claims = jwt.MapClaims{
		"exp": expirationTime,
	}

	secretKey := []byte("SUPER_SECRET_KEY")
	return token.SignedString(secretKey)
}
