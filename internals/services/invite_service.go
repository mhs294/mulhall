package services

import (
	"fmt"
	"time"

	"github.com/mhs294/mulhall/internals/db"
	"github.com/mhs294/mulhall/internals/env"
	"github.com/mhs294/mulhall/internals/types"
	"github.com/mhs294/mulhall/internals/utils"
)

// InviteService represents a service for interacting with User Invites for the site.
type InviteService struct {
	inviteRepository *db.InviteRepository
}

// NewInviteService creates a new instance of an InviteService and returns a pointer to it.
//
// ir is the pointer to the InviteRepository that will be used at runtime by the InviteService.
func NewInviteService(ir *db.InviteRepository) *InviteService {
	return &InviteService{inviteRepository: ir}
}

// CreateInvite creates a new Invite from the provided information and returns a pointer to it.
//
// req is the CreateInviteRequest containing the necessary information to create the new Invite.
func (s *InviteService) CreateInvite(req *types.CreateInviteRequest) (*types.Invite, error) {
	// Create the Invite with a randomly generated validation token
	token := utils.CreateAlphaNumToken(64)
	invite := &types.Invite{
		Email:        req.Email,
		Contestant:   req.Contestant,
		Role:         req.Role,
		InvitingUser: req.InvtingUser,
		Token:        token,
		Expiration:   time.Now().Add(env.InviteExpiration),
		Accepted:     false,
	}

	// Insert the Invite into the database
	if err := s.inviteRepository.InsertInvite(invite); err != nil {
		return nil, fmt.Errorf("failed to create invite: %v", err)
	}

	return invite, nil
}

// ValidateInvite loads the Invite for the proided email/token combination and
// returns the Invite's unique identifier if it is active.
// Returns an error if the Invite does not exist for the provided email/token
// combination, or the Invite has expired.
//
// email is the email address to look up the Invite for.
//
// token is the token string that should match with the email on the Invite.
func (s *InviteService) ValidateInvite(email string, token string) (types.InviteID, error) {
	inv, err := s.inviteRepository.GetInvite(email, token)
	if err != nil {
		return "", err
	}

	if inv == nil {
		return "", &types.InviteNotFoundError{Email: email, Token: token}
	}

	if inv.Expiration.Before(time.Now()) {
		return "", &types.InviteExpiredError{Email: email, Token: token}
	}

	return inv.ID, nil
}

// AcceptInvite marks the Invite for the provided email/token combination as accepted.
//
// Returns an error if the Invite does not exist for the provided email/token
// combination, or the Invite has expired.
//
// email is the email address to look up the Invite for.
//
// token is the token string that should match with the email on the Invite.
func (s *InviteService) AcceptInvite(email string, token string) error {
	id, err := s.ValidateInvite(email, token)
	if err != nil {
		return err
	}

	return s.inviteRepository.AcceptInvite(id)
}
