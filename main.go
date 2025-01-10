package main

import (
	"awesomeProject/Config"
	"awesomeProject/Controller"
	"github.com/labstack/echo/v4"
	"log"
)

func main() {
	db := Config.InitializeDatabase()
	e := echo.New()

	Controller.Router(e, db)

	if err := e.Start(":8080"); err != nil {
		log.Fatal("Error starting the server: ", err)
	}
}
