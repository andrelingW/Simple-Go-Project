package Controller

import (
	"awesomeProject/Config"
	"awesomeProject/Model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

func Router(e *echo.Echo, db *gorm.DB) {
	e.POST("/login", LoginHandler(db)) // Public route
	e.POST("/register", RegisterHandlers(db))

	// Secured route with JWT validation middleware
	e.GET("/protected", Config.Middleware(TestHandler)) // Secured route

}

func RegisterHandlers(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user Model.UserModel

		// Bind JSON request body to the user struct
		if err := c.Bind(&user); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid JSON"})
		}

		// Insert the user into the database using GORM
		if err := db.Create(&user).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create user"})
		}

		// Return the created user as a response
		return c.JSON(http.StatusOK, "Successfully created user")
	}
}

func LoginHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var loginData struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		// Bind the request body to the loginData struct
		if err := c.Bind(&loginData); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request data"})
		}

		// Check if user exists in the database based on email
		var user Model.UserModel
		result := db.Where("email = ?", loginData.Email).First(&user)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "User not found"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
		}

		// For now, we'll assume passwords match. You should hash and compare passwords in a real app.
		if user.Password != loginData.Password {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid credentials"})
		}
		token, err := Config.GenerateJWT()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate JWT")
		}

		return c.JSON(http.StatusOK, map[string]string{
			"token": token,
		})
	}
}

func TestHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "This is a protected route. You have access!",
	})
}
