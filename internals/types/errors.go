package types

import (
	"fmt"
	"time"
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

// The system attempted to find a User that does not exist or has been deactivated.
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
type SessionExpiredError struct {
	ID SessionID
}

func (e *SessionExpiredError) Error() string {
	return fmt.Sprintf("session is expired. id=%s", e.ID)
}

// The system attempted to find a Pool that does not exist or has been deactivated.
type PoolNotFoundError struct {
	ID PoolID
}

func (e *PoolNotFoundError) Error() string {
	return fmt.Sprintf("failed to find pool. id=%s", e.ID)
}

// The system attempted to find a Contestant that does not exist or has been deactivated.
type ContestantNotFoundError struct {
	ID ContestantID
}

func (e *ContestantNotFoundError) Error() string {
	return fmt.Sprintf("failed to find contestant. id=%s", e.ID)
}

// The system attempted to find a Schedule that does not exist or has been deactivated.
type ScheduleNotFoundError struct {
	Date time.Time
	Year int
	Week int
}

func (e *ScheduleNotFoundError) Error() string {
	if e.Date.IsZero() {
		return fmt.Sprintf("failed to find schedule. date=%s", e.Date.Format(time.UnixDate))
	}

	return fmt.Sprintf("failed to find schedule. year=%d week=%d.", e.Year, e.Week)
}

// The system attempted to create a new Schedule for a period which already has an existing Schedule.
type ScheduleConflictError struct {
	Date string
	Year int
	Week int
	ID   ScheduleID
}

func (e *ScheduleConflictError) Error() string {
	if len(e.Date) > 0 {
		return fmt.Sprintf("a schedule already exists for date=%s (id=%v).", e.Date, e.ID)
	}

	return fmt.Sprintf("a schedule already exists for year=%d/week=%d (id=%v).", e.Year, e.Week, e.ID)
}

// The system attempted to create a new Schedule with a closing date/time that falls outside the start and end of the week.
type ScheduleInvalidClosesError struct {
	Start   time.Time
	End     time.Time
	Request *CreateScheduleRequest
}

func (e *ScheduleInvalidClosesError) Error() string {
	return fmt.Sprintf("schedule closes date/time must fall between %s and %s (request=%v).",
		e.Start.Format(time.UnixDate),
		e.End.Format(time.UnixDate),
		e.Request)
}
