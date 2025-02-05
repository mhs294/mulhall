package ioc

import (
	"github.com/mhs294/mulhall/internals/services"
)

var invService *services.InviteService
var userService *services.UserService
var poolService *services.PoolService
var conService *services.ContestantService

func InviteService() *services.InviteService {
	if invService == nil {
		repo := InviteRepository()
		invService = services.NewInviteService(repo)
	}

	return invService
}

func UserService() *services.UserService {
	if userService == nil {
		invServ := InviteService()
		userRepo := UserRepository()
		sessRepo := SessionRepository()
		userService = services.NewUserService(invServ, userRepo, sessRepo)
	}

	return userService
}

func PoolService() *services.PoolService {
	if poolService == nil {
		repo := PoolRepository()
		poolService = services.NewPoolService(repo)
	}

	return poolService
}

func ContestantService() *services.ContestantService {
	if conService == nil {
		conRepo := ContestantRepository()
		poolService := PoolService()
		conService = services.NewContestantService(conRepo, poolService)
	}

	return conService
}
