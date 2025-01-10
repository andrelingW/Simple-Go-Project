package Config

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get token from Authorization header
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization Header")
		}

		// Bearer <token>
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Authorization Header")
		}

		// Parse the JWT token
		secretKey := []byte("SUPER_SECRET_KEY")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Make sure the token method is what we expect
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid Signing Method")
			}
			return secretKey, nil
		})

		// Handle parsing errors
		if err != nil || !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or Expired Token")
		}

		// Continue to the next middleware or handler if the token is valid
		return next(c)
	}
}
