package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

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
		acc.POST("/register", c.register)
		acc.POST("/login", c.login)
	}
}

func (c *UserController) register(ctx *gin.Context) {
	var req *types.RegisterUserRequest
	if err := utils.FromRequestJSON(&req, ctx); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		c.logger.Printf("failed to unmarhsal RegisterUserRequest from json: %v", err)
		return
	}

	if _, err := c.userService.Register(req); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *UserController) login(ctx *gin.Context) {
	var req *types.LoginRequest
	if err := utils.FromRequestJSON(&req, ctx); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		c.logger.Printf("failed to unmarhsal LoginRequest from json: %v", err)
		return
	}

	sess, err := c.userService.Login(req.Email, req.Password)
	if err != nil {
		switch err.(type) {
		case *types.UserNotFoundError:
		case *types.PasswordIncorrectError:
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		default:
			ctx.AbortWithStatus(http.StatusInternalServerError)
			c.logger.Printf("unexpected error occurred while attempting to login: %v", err)
			return
		}
	}

	maxAge := int64(sess.Expiration.Sub(time.Now().UTC()).Seconds())
	cookie := fmt.Sprintf("mulhall.sessionID=%s; Max-Age=%d", sess.ID, maxAge)
	ctx.Header(http.CanonicalHeaderKey("set-cookie"), cookie)
	ctx.Status(http.StatusOK)
}
