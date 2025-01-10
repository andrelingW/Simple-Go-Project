package main

import (
	"awesomeProject/Config"
	"awesomeProject/Controller"
	"github.com/labstack/echo/v4"
	"log"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	// Initialize the SQLite database connection
	db := Config.InitializeDatabase()

	// Initialize the Echo router
	e := echo.New()

	// Set up routes
	Controller.Router(e, db)

	// Start the server
	if err := e.Start(":8080"); err != nil {
		log.Fatal("Error starting the server: ", err)
	}
}
