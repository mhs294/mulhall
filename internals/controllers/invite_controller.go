package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mhs294/mulhall/internals/services"
	"github.com/mhs294/mulhall/internals/types"
	"github.com/mhs294/mulhall/internals/utils"
)

// InviteController is responsible for handling requests for Invite HTTP APIs.
type InviteController struct {
	logger        *log.Logger
	inviteService *services.InviteService
}

// NewInviteController creates a new instance of an InviteController and returns a pointer to it.
//
// l is the pointer to the [log.Logger] that will be used at runtime by the InviteController.
//
// s is the pointer to the InviteService that will be used at runtime by the InviteController.
func NewInviteController(l *log.Logger, s *services.InviteService) *InviteController {
	return &InviteController{logger: l, inviteService: s}
}

// RegisterHandlers defines this controller's HTTP routes and their corresponding handler functions.
func (c *InviteController) RegisterHandlers(e *gin.Engine) {
	inv := e.Group("/invite")
	{
		inv.POST("/create", c.create)
		inv.POST("/accept", c.accept)
	}
}

func (c *InviteController) create(ctx *gin.Context) {
	var req *types.CreateInviteRequest
	if err := utils.FromRequestJSON(&req, ctx); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		c.logger.Printf("failed to unmarhsal CreateInviteRequest from json: %v", err)
		return
	}

	_, err := c.inviteService.CreateInvite(req)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		c.logger.Printf("failed to create invite: %v", err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *InviteController) accept(ctx *gin.Context) {
	email := ctx.Query("email")
	if len(email) == 0 {
		c.logger.Printf("attempted to validate an invite with no email")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	token := ctx.Query("token")
	if len(token) == 0 {
		c.logger.Printf("attempted to validate an invite with no token")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if _, err := c.inviteService.ValidateInvite(email, token); err != nil {
		switch err.(type) {
		case *types.InviteNotFoundError:
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		case *types.InviteExpiredError:
			ctx.AbortWithStatus(http.StatusGone)
			return
		case *types.InviteAlreadyAcceptedError:
			ctx.AbortWithStatus(http.StatusConflict)
			return
		default:
			ctx.AbortWithStatus(http.StatusInternalServerError)
			c.logger.Printf("unexpected error occurred while validating invite: %v", err)
			return
		}
	}

	ctx.Status(http.StatusNoContent)
}
