package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/mhs294/mulhall/internals/db"
	"github.com/mhs294/mulhall/internals/env"
	"github.com/mhs294/mulhall/views"
)

// ViewController is responsible for serving HTML views to the end user via HTTP.
type ViewController struct {
	TeamRepo *db.TeamRepository
}

// NewViewController creates a new instance of a ViewController and returns a pointer to it.
func NewViewController() (*ViewController, error) {
	teamRepo, err := db.NewTeamRepository(env.MongoDBConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create new ViewController: %v", err)
	}

	return &ViewController{TeamRepo: teamRepo}, nil
}

// RegisterHandlers defines this controller's HTTP routes and their corresponding handler functions.
func (c *ViewController) RegisterHandlers(e *gin.Engine) {
	e.GET("/", c.index)
}

func (c *ViewController) index(ctx *gin.Context) {
	_, cancel := context.WithTimeout(context.Background(), env.Timeout)
	defer cancel()

	teams := c.TeamRepo.GetAllTeams()
	render(ctx, http.StatusOK, views.Index(teams))
}

func render(ctx *gin.Context, status int, template templ.Component) error {
	ctx.Status(status)
	return template.Render(ctx.Request.Context(), ctx.Writer)
}
