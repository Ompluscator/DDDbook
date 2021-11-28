package presentation

import (
	"crypto/md5"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"

	"github.com/PacktPublishing/Domain-Driven-Design-with-Go/chapter5/internal/user/application"
	"github.com/PacktPublishing/Domain-Driven-Design-with-Go/chapter5/internal/user/domain"
)

// userJSON is a DTO used for JSON body in REST API
type userJSON struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// newUserJSON creates new DTO from domain.User
func newUserJSON(user domain.User) userJSON {
	return userJSON{
		ID:        user.ID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
}

// registrationJSON is a DTO used for JSON body in REST API
type registrationJSON struct {
	Username        string `json:"username" validate:"required"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
}

// toEntity creates new domain.User from DTO
func (j registrationJSON) toEntity() domain.User {
	return domain.User{
		Username:  j.Username,
		FirstName: j.FirstName,
		LastName:  j.LastName,
	}
}

// UserController represents the main controller for domain.User
type UserController struct {
	registration application.RegistrationUseCase
	validate     *validator.Validate
}

// NewUserController creates new UserController
func NewUserController(registration application.RegistrationUseCase, validate *validator.Validate) *UserController {
	return &UserController{
		registration: registration,
		validate:     validate,
	}
}

// Register intercepts the request for registration and push it to the layer below
func (c *UserController) Register(ctx *gin.Context) {
	var registrationData registrationJSON

	if err := ctx.ShouldBindJSON(&registrationData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.validate.StructCtx(ctx, registrationData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	passwordData := []byte(registrationData.Password)
	hash := fmt.Sprintf("%x", md5.Sum(passwordData))

	result, err := c.registration.Execute(ctx, registrationData.toEntity(), hash)
	if err != nil {
		switch err {
		case application.RegistrationUseCaseUserAlreadyCreated:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, newUserJSON(*result))
}
