package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mhs294/mulhall/internals/controllers"
	"github.com/mhs294/mulhall/internals/middleware"
)

type Controller interface {
	RegisterHandlers(e *gin.Engine)
}

type AppServer struct {
	Router      *gin.Engine
	Controllers []Controller
}

func NewAppServer() (*AppServer, error) {
	r := initRouter()
	conts, err := initControllers()
	if err != nil {
		return nil, fmt.Errorf("failed to create new AppServer: %v", err)
	}

	for _, c := range conts {
		c.RegisterHandlers(r)
	}

	return &AppServer{Router: r, Controllers: conts}, nil
}

func (s *AppServer) Start() {
	// Port must match EXPOSE command in Dockerfile
	s.Router.Run("0.0.0.0:8080")
}

func initRouter() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	// Static content (CSS/JS)
	r.Static("/static", "./static")

	// Middleware
	r.Use(middleware.CORS)
	return r
}

func initControllers() ([]Controller, error) {
	conts := make([]Controller, 0)

	// ViewController
	viewCon, err := controllers.NewViewController()
	if err != nil {
		return nil, err
	}
	conts = append(conts, viewCon)

	return conts, nil
}
