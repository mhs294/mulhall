package types

import (
	"github.com/mhs294/mulhall/internals/types/roles"
)

// CreateInviteRequest contains all of the information necessary to create an Invite for a new User.
type CreateInviteRequest struct {
	Email       string       `json:"email"`
	Contestant  ContestantID `json:"contestantId"`
	Role        roles.Role   `json:"role"`
	InvtingUser UserID       `json:"invitingUserId"`
}
