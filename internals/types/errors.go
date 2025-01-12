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

// The system attempted to validate or accept an Invite that has already been accepteds.
type InviteAlreadyAcceptedError struct {
	Email string
	Token string
}

func (e *InviteAlreadyAcceptedError) Error() string {
	return fmt.Sprintf("invite has already been accepted. email=%s, token=%s", e.Email, e.Token)
}

// The system attempted to find a User that does not exist.
type UserNotFoundError struct {
	Email string
}

func (e *UserNotFoundError) Error() string {
	return fmt.Sprintf("failed to find user. email=%s", e.Email)
}

// The system attempted to register a new User with a mismatched password/confirm.
type PasswordMismatchError struct{}

func (e *PasswordMismatchError) Error() string {
	return "password and confirm password fields do not match."
}

// The system attempted to login a User with an incorrect password.
type PasswordIncorrectError struct{}

func (e *PasswordIncorrectError) Error() string {
	return "password was incorrect."
}
