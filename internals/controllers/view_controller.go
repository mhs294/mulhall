package controllers

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/mhs294/mulhall/internals/env"
	"github.com/mhs294/mulhall/internals/middleware"
	"github.com/mhs294/mulhall/internals/repos"
	"github.com/mhs294/mulhall/views"
)

// ViewController is responsible for serving HTML views to the end user via HTTP.
type ViewController struct {
	userAuth *middleware.UserAuthMiddleware
	teamRepo *repos.TeamRepository
}

// NewViewController creates a new instance of a ViewController and returns a pointer to it.
func NewViewController(ua *middleware.UserAuthMiddleware, tr *repos.TeamRepository) *ViewController {
	return &ViewController{userAuth: ua, teamRepo: tr}
}

// RegisterHandlers defines this controller's HTTP routes and their corresponding handler functions.
func (c *ViewController) RegisterHandlers(e *gin.Engine) {
	e.GET("/", c.userAuth.ViewAuth, c.index)
	e.GET("/login", c.login)
}

func (c *ViewController) index(ctx *gin.Context) {
	_, cancel := context.WithTimeout(context.Background(), env.Timeout)
	defer cancel()

	teams, err := c.teamRepo.GetAllTeams()
	if err != nil {
		// TODO - replace with error view
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	render(ctx, http.StatusOK, views.Index(teams))
}

func (c *ViewController) login(ctx *gin.Context) {
	// TODO - render login view
	render(ctx, http.StatusOK, views.Login())
}

func render(ctx *gin.Context, status int, template templ.Component) {
	if err := template.Render(ctx.Request.Context(), ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	ctx.Status(status)
}
