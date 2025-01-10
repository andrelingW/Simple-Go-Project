package Config

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var jwtSecret = []byte("SUPER_SECRET_KEY") // Replace with your secret key

func GenerateJWT() (string, error) {
	// Create the token with claims and sign it
	token := jwt.New(jwt.SigningMethodHS256)

	expirationTime := time.Now().Add(1 * time.Hour).Unix()
	token.Claims = jwt.MapClaims{
		"exp": expirationTime,
	}

	// Sign the token using a secret key
	secretKey := []byte("SUPER_SECRET_KEY")
	return token.SignedString(secretKey)
}

// ValidateJWT parses and validates a JWT token
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})
}
