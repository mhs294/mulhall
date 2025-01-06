package ioc

import (
	"github.com/mhs294/mulhall/internals/services"
)

var inviteService *services.InviteService
var userService *services.UserService

func InviteService() *services.InviteService {
	if inviteService == nil {
		repo := InviteRepository()
		inviteService = services.NewInviteService(repo)
	}

	return inviteService
}

func UserService() *services.UserService {
	if userService == nil {
		invServ := InviteService()
		repo := UserRepository()
		userService = services.NewUserService(invServ, repo)
	}

	return userService
}
