package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mhs294/mulhall/internals/ioc"
	"github.com/mhs294/mulhall/internals/middleware"
)

// Controller is a type that has HTTP method handlers capable of being registered with the app server engine.
type Controller interface {
	RegisterHandlers(e *gin.Engine)
}

// AppServer represents the backend server responsible for serving views to the end user
// as well as handling HTTP API requests to facilitate user workflows in the web application.
type AppServer struct {
	Router *gin.Engine
}

// NewAppServer constructs a new instance of an AppServer and returns a pointer to it.
func NewAppServer() (*AppServer, error) {
	r := initRouter()

	conts := initControllers()
	for _, c := range conts {
		c.RegisterHandlers(r)
	}

	return &AppServer{Router: r}, nil
}

// Start runs the application server.
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

func initControllers() []Controller {
	conts := make([]Controller, 0)
	conts = append(conts, ioc.InviteController())
	conts = append(conts, ioc.UserController())
	conts = append(conts, ioc.ViewController())

	return conts
}
