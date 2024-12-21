package types

import (
	"fmt"
)

// The system attempted to find an Invite that does not exist.
type InviteNotFoundError struct {
	Email string
	Token string
}

func (e *InviteNotFoundError) Error() string {
	return fmt.Sprintf("failed to find invite. email=%s, token=%s", e.Email, e.Token)
}

// The system attempted to validate or accept an expired Invite.
type InviteExpiredError struct {
	Email string
	Token string
}

func (e *InviteExpiredError) Error() string {
	return fmt.Sprintf("invite expired. email=%s, token=%s", e.Email, e.Token)
}
