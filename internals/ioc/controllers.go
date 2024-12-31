package ioc

import (
	"github.com/mhs294/mulhall/internals/controllers"
)

var inviteCont *controllers.InviteController
var accountCont *controllers.AccountController
var viewCont *controllers.ViewController

func InviteController() *controllers.InviteController {
	if inviteCont == nil {
		logger := Logger()
		service := InviteService()
		inviteCont = controllers.NewInviteController(logger, service)
	}

	return inviteCont
}

func AccountController() *controllers.AccountController {
	if accountCont == nil {
		logger := Logger()
		service := AccountService()
		accountCont = controllers.NewAccountController(logger, service)
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
