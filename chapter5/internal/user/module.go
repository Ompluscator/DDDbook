package user

import (
	"log"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/PacktPublishing/Domain-Driven-Design-with-Go/chapter5/internal/user/application"
	"github.com/PacktPublishing/Domain-Driven-Design-with-Go/chapter5/internal/user/infrastructure"
	"github.com/PacktPublishing/Domain-Driven-Design-with-Go/chapter5/internal/user/presentation"
)

// Module is a struct that defines all dependencies inside user module
type Module struct{}

// Configure setups all dependencies
func (m *Module) Configure(databasePath string, engine *gin.Engine, validate *validator.Validate) {
	db, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	repository := infrastructure.NewUserRepository(db)
	useCase := application.NewRegistrationUseCase(repository)
	controller := presentation.NewUserController(useCase, validate)
	engine.POST("/users", controller.Register)
}
