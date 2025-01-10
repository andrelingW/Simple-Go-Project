package Controller

import (
	"awesomeProject/Config"
	"github.com/labstack/echo/v4"
	"net/http"
)

func LoginHandler(c echo.Context) error {
	token, err := Config.GenerateJWT()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate JWT")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

func TestHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "This is a protected route. You have access!",
	})
}
