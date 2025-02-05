package types

import (
	"time"

	"github.com/mhs294/mulhall/internals/types/roles"
	"github.com/mhs294/mulhall/internals/types/status"
)

// Team contains the basic identifying details of an NFL team.
type Team struct {
	ID        TeamID `json:"id"`
	Shorthand string `json:"shorthand"`
	Location  string `json:"location"`
	Name      string `json:"name"`
}

// User represents an individual person within the site and their account details.
type User struct {
	ID            UserID `json:"id"`
	Email         string `json:"email"`
	Salt          string `json:"-"` // Always omit this field from JSON serialization
	Hash          string `json:"-"` // Always omit this field from JSON serialization
	Administrator bool   `json:"administrator"`
	Active        bool   `json:"active"`
}

// Invite represents an invitation for a new user to create an account with the site and join a Contestant.
type Invite struct {
	ID           InviteID     `json:"id"`
	Email        string       `json:"email"`
	Contestant   ContestantID `json:"contestant"`
	Role         roles.Role   `json:"role"`
	InvitingUser UserID       `json:"invitingUser"`
	Token        string       `json:"-"` // Always omit this field from JSON serialization
	Expiration   time.Time    `json:"expiration"`
	Accepted     bool         `json:"accepted"`
}

// Session represents an authentication session for a logged in user.
type Session struct {
	ID         SessionID
	User       UserID
	Expiration time.Time
}

// Pool defines a set of rules for an elimination game in which a group of Contestants compete.
type Pool struct {
	ID          PoolID                    `json:"id"`
	Name        string                    `json:"name"`
	Contestants map[ContestantID]struct{} `json:"contestants"`
	Active      bool                      `json:"active"`
	Complete    bool                      `json:"complete"`
}

// Contestant defines a single entry within the pool and the authorized Users who maintain it.
type Contestant struct {
	ID              ContestantID          `json:"id"`
	Name            string                `json:"name"`
	AuthorizedUsers map[UserID]roles.Role `json:"authorizedUsers"`
	Status          status.Status         `json:"status"`
	Active          bool                  `json:"active"`
}

// Schedule represents a set of Matchups that occur within a given week.
type Schedule struct {
	ID       ScheduleID  `json:"id"`
	Year     int         `json:"year"`
	Week     int         `json:"week"`
	Matchups []MatchupID `json:"matchups"`
}

// Matchup represents an individual game between two Teams.
type Matchup struct {
	ID       MatchupID `json:"id"`
	AwayTeam TeamID    `json:"awayTeam"`
	HomeTeam TeamID    `json:"homeTeam"`
	DateTime time.Time `json:"dateTime"`
}

// Entry represents a Contestant's Pick for a given Schedule, as well as any potential suggested Picks.
type Entry struct {
	ID             EntryID    `json:"id"`
	Schedule       ScheduleID `json:"schedule"`
	SelectedPick   Pick       `json:"selectedPick"`
	SuggestedPicks []Pick     `json:"suggestedPicks"`
}

// Pick represents a selection of a Team within a given Matchup.
type Pick struct {
	Matchup MatchupID `json:"matchup"`
	Team    TeamID    `json:"team"`
}
