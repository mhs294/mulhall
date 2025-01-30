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
	ID    UserID
	Email string
}

func (e *UserNotFoundError) Error() string {
	detail := ""
	if len(e.ID) > 0 {
		detail = fmt.Sprintf("id=%s", e.ID)
	} else if len(e.Email) > 0 {
		detail = fmt.Sprintf("email=%s", e.Email)
	}

	return fmt.Sprintf("failed to find user. %s", detail)
}

// The system attempted to load a User that is inactive.
type UserInactiveError struct {
	ID    UserID
	Email string
}

func (e *UserInactiveError) Error() string {
	detail := ""
	if len(e.ID) > 0 {
		detail = fmt.Sprintf("id=%s", e.ID)
	} else if len(e.Email) > 0 {
		detail = fmt.Sprintf("email=%s", e.Email)
	}

	return fmt.Sprintf("user is inactive. %s", detail)
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

// The system attempted to authenticate a request with a missing/empty Session-ID header.
type MissingSessionIDError struct{}

func (e *MissingSessionIDError) Error() string {
	return "Session-ID request header was missing or blank."
}

// The system attempted to find a Session that does not exist.
type SessionNotFoundError struct {
	ID SessionID
}

func (e *SessionNotFoundError) Error() string {
	return fmt.Sprintf("failed to find session. id=%s", e.ID)
}

// The system attempted to authenticate a request with an expired Session.
type SessionExpiredError struct{}

func (e *SessionExpiredError) Error() string {
	return "session is expired."
}
