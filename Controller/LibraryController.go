package Controller

import (
	"awesomeProject/Config"
	"awesomeProject/Model"
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

func Router(e *echo.Echo, db *gorm.DB) {
	//Public
	e.POST("/login", loginHandler(db))
	e.POST("/register", registerHandlers(db))

	//Secured
	e.GET("/view/books", Config.Middleware(viewAllBookHandler(db))) // Secured route
	e.GET("/view/description/:book", Config.Middleware(viewBookDetailHandler(db)))
	e.GET("/view/borrow/:id", Config.Middleware(borrowBookHandler(db))) // Secured route
	e.GET("/view/return/:id", Config.Middleware(returnBookHandler(db)))
}

func registerHandlers(db *gorm.DB) echo.HandlerFunc {
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

func loginHandler(db *gorm.DB) echo.HandlerFunc {
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

func viewAllBookHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var books []Model.BookModel

		type BookResponse struct {
			ID        int    `json:"id"`
			Title     string `json:"title"`
			Author    string `json:"author"`
			Available bool   `json:"available"`
		}

		if err := db.Find(&books).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve books"})
		}

		// Transform the result into the response struct
		var response []BookResponse
		for _, book := range books {
			response = append(response, BookResponse{
				ID:        book.ID,
				Title:     book.Title,
				Author:    book.Author,
				Available: book.Available,
			})
		}

		return c.JSON(http.StatusOK, response)
	}
}

func viewBookDetailHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.Param("book")
		var book Model.BookModel

		// Query the database for a book with the specified ID
		if err := db.First(&book, name).Error; err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "Book not found"})
		}

		return c.JSON(http.StatusOK, book)
	}
}

func borrowBookHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the book ID from the request parameters
		bookID := c.Param("id")

		// Find the book by ID
		var book Model.BookModel
		if err := db.First(&book, bookID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.JSON(http.StatusNotFound, map[string]string{"message": "Book not found"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to find book"})
		}

		// Check if the book is already borrowed
		if !book.Available {
			return c.JSON(http.StatusConflict, map[string]string{"message": "Book is already borrowed"})
		}

		// Set `Available` to `false`
		book.Available = false
		if err := db.Save(&book).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to borrow book"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Book borrowed successfully"})
	}
}

func returnBookHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the book ID from the request parameters
		bookID := c.Param("id")

		// Find the book by ID
		var book Model.BookModel
		if err := db.First(&book, bookID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.JSON(http.StatusNotFound, map[string]string{"message": "Book not found"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to find book"})
		}

		// Check if the book is already returned
		if book.Available {
			return c.JSON(http.StatusConflict, map[string]string{"message": "Book is already returned"})
		}

		// Set `Available` to `true`
		book.Available = true
		if err := db.Save(&book).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to return book"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Book returned successfully"})
	}
}
