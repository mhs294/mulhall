package ioc

import (
	"github.com/mhs294/mulhall/internals/middleware"
)

var userAuthMiddleWare *middleware.UserAuthMiddleware

func UserAuthMiddleware() *middleware.UserAuthMiddleware {
	if userAuthMiddleWare == nil {
		logger := Logger()
		sessRepo := SessionRepository()
		userAuthMiddleWare = middleware.NewUserAuthMiddleware(logger, sessRepo)
	}

	return userAuthMiddleWare
}
