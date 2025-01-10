package Config

import (
	"awesomeProject/Model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"log"
)

// Initialize initializes the SQLite database and returns a pointer to the database
func InitializeDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Model.UserModel{})

	log.Println("Database migration successful")

	return db
}
