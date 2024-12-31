package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mhs294/mulhall/internals/services"
	"github.com/mhs294/mulhall/internals/types"
	"github.com/mhs294/mulhall/internals/utils"
)

// AccountController is responsible for handling requests for Account HTTP APIs.
type AccountController struct {
	logger         *log.Logger
	accountService *services.AccountService
}

// NewAccountController creates a new instance of an AccountController and returns a pointer to it.
//
// l is the pointer to the [log.Logger] that will be used at runtime by the AccountController.
//
// s is the pointer to the AccountService that will be used at runtime by the AccountController.
func NewAccountController(l *log.Logger, s *services.AccountService) *AccountController {
	return &AccountController{logger: l, accountService: s}
}

// RegisterHandlers defines this controller's HTTP routes and their corresponding handler functions.
func (c *AccountController) RegisterHandlers(e *gin.Engine) {
	acc := e.Group("/account")
	{
		acc.POST("/register", c.registerAccount)
	}
}

func (c *AccountController) registerAccount(ctx *gin.Context) {
	var req types.RegisterAccountRequest
	if err := utils.FromRequestJSON(&req, ctx); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO - replace with call to account service to validate and save new account
	ctx.JSON(http.StatusOK, req)
}
