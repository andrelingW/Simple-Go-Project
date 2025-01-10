package Controller

import (
	"awesomeProject/Config"
	"awesomeProject/Model"
	_ "awesomeProject/docs"
	"errors"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"
	"net/http"
)

func Router(e *echo.Echo, db *gorm.DB) {
	// Public
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.POST("/login", loginHandler(db))
	e.POST("/register", registerHandlers(db))

	// Secured
	e.GET("/view/books", Config.Middleware(viewAllBookHandler(db)))
	e.GET("/view/description/:book", Config.Middleware(viewBookDetailHandler(db)))
	e.GET("/view/borrow/:id", Config.Middleware(borrowBookHandler(db)))
	e.GET("/view/return/:id", Config.Middleware(returnBookHandler(db)))
}

// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags users
// @Accept json
// @Produce json
// @Param user body Model.UserModel true "User registration details"
// @Success 200 {string} string "Successfully created user"
// @Failure 400 {object} map[string]string "Invalid JSON"
// @Failure 500 {object} map[string]string "Failed to create user"
// @Router /register [post]
func registerHandlers(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user Model.UserModel

		if err := c.Bind(&user); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid JSON"})
		}

		if err := db.Create(&user).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create user"})
		}

		return c.JSON(http.StatusOK, "Successfully created user")
	}
}

// @Summary User login
// @Description Login user and receive a JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body map[string]string true "Login credentials"
// @Success 200 {object} map[string]string "JWT token"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /login [post]
func loginHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var loginData struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.Bind(&loginData); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request data"})
		}

		var user Model.UserModel
		result := db.Where("email = ?", loginData.Email).First(&user)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "User not found"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
		}

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

// @Summary Get all books
// @Description Retrieve all books excluding their descriptions
// @Tags books
// @Produce json
// @Success 200 {array} Model.BookResponse "List of books"
// @Failure 500 {object} map[string]string "Failed to retrieve books"
// @Router /view/books [get]
func viewAllBookHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var books []Model.BookModel

		if err := db.Find(&books).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve books"})
		}

		var response []Model.BookResponse
		for _, book := range books {
			response = append(response, Model.BookResponse{
				ID:        book.ID,
				Title:     book.Title,
				Author:    book.Author,
				Available: book.Available,
			})
		}

		return c.JSON(http.StatusOK, response)
	}
}

// @Summary Get book details
// @Description Retrieve detailed information about a specific book
// @Tags books
// @Produce json
// @Param book path string true "Book ID"
// @Success 200 {object} Model.BookModel "Book details"
// @Failure 404 {object} map[string]string "Book not found"
// @Router /view/description/{book} [get]
func viewBookDetailHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		bookID := c.Param("id")
		var book Model.BookModel

		if err := db.First(&book, bookID).Error; err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "Book not found"})
		}

		return c.JSON(http.StatusOK, book)
	}
}

// @Summary Borrow a book
// @Description Mark a book as borrowed
// @Tags books
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} map[string]string "Book borrowed successfully"
// @Failure 404 {object} map[string]string "Book not found"
// @Failure 409 {object} map[string]string "Book is already borrowed"
// @Failure 500 {object} map[string]string "Failed to borrow book"
// @Router /view/borrow/{id} [get]
func borrowBookHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		bookID := c.Param("id")
		var book Model.BookModel

		if err := db.First(&book, bookID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.JSON(http.StatusNotFound, map[string]string{"message": "Book not found"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to find book"})
		}

		if !book.Available {
			return c.JSON(http.StatusConflict, map[string]string{"message": "Book is already borrowed"})
		}

		book.Available = false
		if err := db.Save(&book).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to borrow book"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Book borrowed successfully"})
	}
}

// @Summary Return a book
// @Description Mark a book as returned
// @Tags books
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} map[string]string "Book returned successfully"
// @Failure 404 {object} map[string]string "Book not found"
// @Failure 409 {object} map[string]string "Book is already returned"
// @Failure 500 {object} map[string]string "Failed to return book"
// @Router /view/return/{id} [get]
func returnBookHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		bookID := c.Param("id")
		var book Model.BookModel

		if err := db.First(&book, bookID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.JSON(http.StatusNotFound, map[string]string{"message": "Book not found"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to find book"})
		}

		if book.Available {
			return c.JSON(http.StatusConflict, map[string]string{"message": "Book is already returned"})
		}

		book.Available = true
		if err := db.Save(&book).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to return book"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Book returned successfully"})
	}
}
