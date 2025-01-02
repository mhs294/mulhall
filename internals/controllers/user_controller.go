package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mhs294/mulhall/internals/services"
	"github.com/mhs294/mulhall/internals/types"
	"github.com/mhs294/mulhall/internals/utils"
)

// UserController is responsible for handling requests for Account HTTP APIs.
type UserController struct {
	logger      *log.Logger
	userService *services.UserService
}

// NewUserController creates a new instance of a UserController and returns a pointer to it.
//
// l is the pointer to the [log.Logger] that will be used at runtime by the UserController.
//
// s is the pointer to the UserService that will be used at runtime by the UserController.
func NewUserController(l *log.Logger, s *services.UserService) *UserController {
	return &UserController{logger: l, userService: s}
}

// RegisterHandlers defines this controller's HTTP routes and their corresponding handler functions.
func (c *UserController) RegisterHandlers(e *gin.Engine) {
	acc := e.Group("/user")
	{
		acc.POST("/register", c.registerUser)
	}
}

func (c *UserController) registerUser(ctx *gin.Context) {
	var req types.RegisterUserRequest
	if err := utils.FromRequestJSON(&req, ctx); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO - replace with call to user service to validate and create new user
	ctx.JSON(http.StatusOK, req)
}
