package main

import (
	"awesomeProject/Config"
	"awesomeProject/Controller"
	"github.com/labstack/echo/v4"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	e := echo.New()

	// Apply JWT middleware to all routes except "/login"
	e.GET("/login", Controller.LoginHandler) // Public route

	// Secured route with JWT validation middleware
	e.GET("/protected", Config.JWTMiddleware(Controller.TestHandler)) // Secured route

	// Start the server
	e.Logger.Fatal(e.Start(":8080"))

}
