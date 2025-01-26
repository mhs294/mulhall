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

// Contestant defines a single entry within the pool and the AuthorizedUsers who maintain it.
type Contestant struct {
	ID              ContestantID     `json:"id"`
	Name            string           `json:"name"`
	AuthorizedUsers []AuthorizedUser `json:"authorizedUsers"`
	Status          status.Status    `json:"status"`
}

// AuthorizedUser is a composite of a User's ID and their corresponding Role for the associated Contestant.
type AuthorizedUser struct {
	ID   UserID     `json:"id"`
	Role roles.Role `json:"role"`
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
	UserID     UserID
	Expiration time.Time
}
