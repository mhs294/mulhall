package ioc

import (
	"github.com/mhs294/mulhall/internals/repos"
)

var teamRepo *repos.TeamRepository
var inviteRepo *repos.InviteRepository
var userRepo *repos.UserRepository
var sessionRepo *repos.SessionRepository
var poolRepo *repos.PoolRepository
var contestantRepo *repos.ContestantRepository
var entryRepo *repos.EntryRepository
var scheduleRepo *repos.ScheduleRepository
var matchupRepo *repos.MatchupRepository

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

func PoolRepository() *repos.PoolRepository {
	if poolRepo == nil {
		mdb := MongoDB()
		poolRepo = repos.NewPoolRepository(mdb)
		if err := poolRepo.TestConnection(); err != nil {
			panic(err)
		}
	}

	return poolRepo
}

func ContestantRepository() *repos.ContestantRepository {
	if contestantRepo == nil {
		mdb := MongoDB()
		contestantRepo = repos.NewContestantRepository(mdb)
		if err := contestantRepo.TestConnection(); err != nil {
			panic(err)
		}
	}

	return contestantRepo
}

func EntryRepository() *repos.EntryRepository {
	if entryRepo == nil {
		mdb := MongoDB()
		entryRepo = repos.NewEntryRepository(mdb)
		if err := entryRepo.TestConnection(); err != nil {
			panic(err)
		}
	}

	return entryRepo
}

func ScheduleRepository() *repos.ScheduleRepository {
	if scheduleRepo == nil {
		mdb := MongoDB()
		scheduleRepo = repos.NewScheduleRepository(mdb)
		if err := scheduleRepo.TestConnection(); err != nil {
			panic(err)
		}
	}

	return scheduleRepo
}

func MatchupRepository() *repos.MatchupRepository {
	if matchupRepo == nil {
		mdb := MongoDB()
		matchupRepo = repos.NewMatchupRepository(mdb)
		if err := matchupRepo.TestConnection(); err != nil {
			panic(err)
		}
	}

	return matchupRepo
}
