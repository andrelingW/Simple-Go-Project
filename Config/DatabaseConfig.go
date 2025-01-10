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
	db.AutoMigrate(&Model.BookModel{})

	log.Println("Database successfully initialized")
	insertBooks(db)

	return db
}

// InsertBooks inserts predefined books into the database with Available = false
func insertBooks(db *gorm.DB) {
	books := []Model.BookModel{
		{Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Description: "A novel about the American Dream.", Available: true},
		{Title: "1984", Author: "George Orwell", Description: "A dystopian novel about totalitarianism.", Available: true},
		{Title: "To Kill a Mockingbird", Author: "Harper Lee", Description: "A novel about racial injustice in the South.", Available: true},
		{Title: "Pride and Prejudice", Author: "Jane Austen", Description: "A story about love, reputation, and class.", Available: true},
		{Title: "The Catcher in the Rye", Author: "J.D. Salinger", Description: "A novel about teenage rebellion and alienation.", Available: true},
		{Title: "Moby-Dick", Author: "Herman Melville", Description: "A tale of obsession and the quest for revenge.", Available: true},
		{Title: "War and Peace", Author: "Leo Tolstoy", Description: "A novel about the Napoleonic wars and Russian society.", Available: true},
		{Title: "The Hobbit", Author: "J.R.R. Tolkien", Description: "A fantasy novel about the journey of Bilbo Baggins.", Available: true},
		{Title: "The Odyssey", Author: "Homer", Description: "An epic poem about Odysseus's journey home.", Available: true},
		{Title: "Crime and Punishment", Author: "Fyodor Dostoevsky", Description: "A psychological drama about guilt and redemption.", Available: true},
	}

	// Insert books into the database
	for _, book := range books {
		if err := db.Create(&book).Error; err != nil {
			log.Println("Error inserting book:", err)
		}
	}

	log.Println("Books have been inserted successfully!")
}
