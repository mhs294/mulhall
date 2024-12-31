package ioc

import (
	"github.com/mhs294/mulhall/internals/repos"
)

var teamRepo *repos.TeamRepository
var inviteRepo *repos.InviteRepository

func TeamRepository() *repos.TeamRepository {
	if teamRepo == nil {
		db := MongoDB()
		teamRepo = repos.NewTeamRepository(db)
		if err := teamRepo.TestConnection(); err != nil {
			panic(err)
		}
	}

	return teamRepo
}

func InviteRepository() *repos.InviteRepository {
	if inviteRepo == nil {
		db := MongoDB()
		inviteRepo = repos.NewInviteRepository(db)
		if err := inviteRepo.TestConnection(); err != nil {
			panic(err)
		}
	}

	return inviteRepo
}
