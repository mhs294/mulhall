package ioc

import (
	"github.com/mhs294/mulhall/internals/controllers"
)

var inviteCont *controllers.InviteController
var accountCont *controllers.UserController
var viewCont *controllers.ViewController

func InviteController() *controllers.InviteController {
	if inviteCont == nil {
		logger := Logger()
		service := InviteService()
		inviteCont = controllers.NewInviteController(logger, service)
	}

	return inviteCont
}

func UserController() *controllers.UserController {
	if accountCont == nil {
		logger := Logger()
		service := UserService()
		accountCont = controllers.NewUserController(logger, service)
	}

	return accountCont
}

func ViewController() *controllers.ViewController {
	if viewCont == nil {
		teamRepo := TeamRepository()
		viewCont = controllers.NewViewController(teamRepo)
	}

	return viewCont
}
