package types

import (
	"time"

	"github.com/mhs294/mulhall/internals/types/roles"
)

// CreateInviteRequest contains all of the information necessary to create an Invite for a new User.
type CreateInviteRequest struct {
	Email          string       `json:"email"`
	ContestantID   ContestantID `json:"contestantId"`
	Role           roles.Role   `json:"role"`
	InvitingUserID UserID       `json:"invitingUserId"`
}

// RegisterUserRequest contains all of the information necessary to register an account for a new User.
type RegisterUserRequest struct {
	Email    string `json:"email"`
	Token    string `json:"token"`
	Password string `json:"password"`
	Confirm  string `json:"confirm"`
}

// RegisterUserRequest contains all of the information necessary to log in a User.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreatePoolRequest contains all of the information necessary to create a new Pool.
type CreatePoolRequest struct {
	Name string `json:"name"`
}

// CreateContestantRequest contains all of the information necessary to create a new Contestant.
type CreateContestantRequest struct {
	Name            string                `json:"name"`
	AuthorizedUsers map[UserID]roles.Role `json:"authorizedUsers"`
}

// CreateScheduleRequest contains all of the information necessary to create a new, empty Schedule.
type CreateScheduleRequest struct {
	Year   int       `json:"year"`
	Week   int       `json:"week"`
	Date   string    `json:"date"` // ISO-8601 date-only string representation (e.g. - "2025-02-09")
	Closes time.Time `json:"closes"`
}

// CreateMatchupRequest contains all of the information necessary to create a new Machup to add to a Schedule.
type CreateMatchupRequest struct {
	ScheduleID ScheduleID `json:"scheduleId"`
	Matchup    *Matchup   `json:"matchup"`
}

// UpdateMatchupRequest contains all of the information necessary to update an existing Matchup within a Schedule.
type UpdateMatchupRequest struct {
	ScheduleID ScheduleID `json:"scheduleId"`
	MatchupID  MatchupID  `json:"matchupId"`
	Matchup    *Matchup   `json:"matchup"`
}

// CreateEntryRequest contains all of the information necessary to create a new Entry for a Contestant in a Pool.
type CreateEntryRequest struct {
	ContestantID ContestantID `json:"contestantId"`
	ScheduleID   ScheduleID   `json:"scheduleId"`
}

// SavePickRequest contains all of the information necessary to save a Selected/Suggested pick for an Entry.
type SavePickRequest struct {
	EntryID   EntryID   `json:"entryId"`
	MatchupID MatchupID `json:"matchupId"`
	TeamID    TeamID    `json:"teamId"`
}
