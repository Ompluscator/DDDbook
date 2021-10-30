package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// User is a struct for handling data from API and database
type User struct {
	ID              uint   `json:"id" gorm:"primaryKey;column:id"`
	Username        string `json:"username" validate:"required" gorm:"column:username"`
	Password        string `json:"password" validate:"required" gorm:"column:password"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password" gorm:"-"`
	FirstName       string `json:"firstName" gorm:"column:firstname"`
	LastName        string `json:"lastName" gorm:"column:lastname"`
}

func main() {
	// setting up a new gin.Engine
	engine := gin.Default()

	// setting up a new connection to the database
	db, err := gorm.Open(sqlite.Open("./database.sqlite"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// settings up a new validator.Validate
	validate := validator.New()

	// defining a new endpoint for user registration
	engine.POST("/users", func(c *gin.Context) {
		var user User

		// mapping the JSON body into User struct
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// validating the user instance
		if err := validate.StructCtx(c, user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// storing the user instance into the database
		if err := db.WithContext(c).Create(&user).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// printing the final created user
		c.JSON(200, user)
	})

	// running the web application
	log.Fatal(engine.Run())
}
