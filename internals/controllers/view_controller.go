package controllers

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/mhs294/mulhall/internals/env"
	"github.com/mhs294/mulhall/internals/repos"
	"github.com/mhs294/mulhall/views"
)

// ViewController is responsible for serving HTML views to the end user via HTTP.
type ViewController struct {
	TeamRepo *repos.TeamRepository
}

// NewViewController creates a new instance of a ViewController and returns a pointer to it.
func NewViewController(tr *repos.TeamRepository) *ViewController {
	return &ViewController{TeamRepo: tr}
}

// RegisterHandlers defines this controller's HTTP routes and their corresponding handler functions.
func (c *ViewController) RegisterHandlers(e *gin.Engine) {
	e.GET("/", c.index)
	e.GET("/login", c.login)
}

func (c *ViewController) index(ctx *gin.Context) {
	_, cancel := context.WithTimeout(context.Background(), env.Timeout)
	defer cancel()

	teams, err := c.TeamRepo.GetAllTeams()
	if err != nil {
		// TODO - replace with error view
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	render(ctx, http.StatusOK, views.Index(teams))
}

func (c *ViewController) login(ctx *gin.Context) {
	// TODO - render login view
}

func render(ctx *gin.Context, status int, template templ.Component) error {
	ctx.Status(status)
	return template.Render(ctx.Request.Context(), ctx.Writer)
}
