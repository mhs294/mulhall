package ioc

import (
	"github.com/mhs294/mulhall/internals/repos"
)

var teamRepo *repos.TeamRepository
var inviteRepo *repos.InviteRepository
var userRepo *repos.UserRepository
var sessionRepo *repos.SessionRepository

func TeamRepository() *repos.TeamRepository {
	if teamRepo == nil {
		mdb := MongoDB()
		teamRepo = repos.NewTeamRepository(mdb)
		if err := teamRepo.TestConnection(); err != nil {
			panic(err)
		}
	}

	return teamRepo
}

func InviteRepository() *repos.InviteRepository {
	if inviteRepo == nil {
		mdb := MongoDB()
		inviteRepo = repos.NewInviteRepository(mdb)
		if err := inviteRepo.TestConnection(); err != nil {
			panic(err)
		}
	}

	return inviteRepo
}

func UserRepository() *repos.UserRepository {
	if userRepo == nil {
		mdb := MongoDB()
		userRepo = repos.NewUserRepository(mdb)
		if err := userRepo.TestConnection(); err != nil {
			panic(err)
		}
	}

	return userRepo
}

func SessionRepository() *repos.SessionRepository {
	if sessionRepo == nil {
		mdb := MongoDB()
		sessionRepo = repos.NewSessionRepository(mdb)
		if err := sessionRepo.TestConnection(); err != nil {
			panic(err)
		}
	}

	return sessionRepo
}
