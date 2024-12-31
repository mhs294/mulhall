package ioc

import (
	"github.com/mhs294/mulhall/internals/services"
)

var inviteService *services.InviteService
var accountService *services.AccountService

func InviteService() *services.InviteService {
	if inviteService == nil {
		repo := InviteRepository()
		inviteService = services.NewInviteService(repo)
	}

	return inviteService
}

func AccountService() *services.AccountService {
	if accountService == nil {
		// TODO - replace this with real service
		accountService = &services.AccountService{}
	}

	return accountService
}
