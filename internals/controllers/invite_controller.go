package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mhs294/mulhall/internals/services"
	"github.com/mhs294/mulhall/internals/types"
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
	invs := e.Group("/invites")
	{
		invs.POST("/validate", c.validateInvite)
		invs.POST("/accept", c.acceptInvite)
	}
}

func (c *InviteController) validateInvite(ctx *gin.Context) {
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

	id, err := c.inviteService.ValidateInvite(email, token)
	if err != nil {
		switch err.(type) {
		case *types.InviteNotFoundError:
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		case *types.InviteExpiredError:
			ctx.AbortWithStatus(http.StatusGone)
			return
		default:
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	ctx.String(http.StatusOK, string(id))
}

func (c *InviteController) acceptInvite(ctx *gin.Context) {
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

	if err := c.inviteService.AcceptInvite(email, token); err != nil {
		switch err.(type) {
		case *types.InviteNotFoundError:
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		case *types.InviteExpiredError:
			ctx.AbortWithStatus(http.StatusGone)
			return
		default:
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	ctx.Status(http.StatusNoContent)
}
