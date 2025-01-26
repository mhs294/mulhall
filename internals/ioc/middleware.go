package ioc

import (
	"github.com/mhs294/mulhall/internals/middleware"
)

var userAuthMiddleWare *middleware.UserAuthMiddleware

func UserAuthMiddleware() *middleware.UserAuthMiddleware {
	if userAuthMiddleWare == nil {
		sessRepo := SessionRepository()
		userAuthMiddleWare = middleware.NewUserAuthMiddleware(sessRepo)
	}

	return userAuthMiddleWare
}
