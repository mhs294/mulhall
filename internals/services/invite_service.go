package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mhs294/mulhall/internals/env"
	"github.com/mhs294/mulhall/internals/repos"
	"github.com/mhs294/mulhall/internals/types"
	"github.com/mhs294/mulhall/internals/utils"
)

// InviteService represents a service for interacting with User Invites for the site.
type InviteService struct {
	invRepo *repos.InviteRepository
}

// NewInviteService creates a new instance of an InviteService and returns a pointer to it.
//
// r is the InviteRepository that will be used at runtime by the InviteService.
func NewInviteService(r *repos.InviteRepository) *InviteService {
	return &InviteService{invRepo: r}
}

// CreateInvite creates a new Invite from the provided information and returns a pointer to it.
//
// req is the CreateInviteRequest containing the necessary information to create the new Invite.
func (s *InviteService) CreateInvite(req *types.CreateInviteRequest) (*types.Invite, error) {
	// Create the Invite with a randomly generated validation token
	token := utils.CreateAlphaNumToken(64)
	inv := &types.Invite{
		ID:           types.InviteID(uuid.New().String()),
		Email:        req.Email,
		Contestant:   req.Contestant,
		Role:         req.Role,
		InvitingUser: req.InvitingUser,
		Token:        token,
		Expiration:   time.Now().UTC().Add(env.InviteExpiration),
		Accepted:     false,
	}

	// Insert the Invite into the database
	if err := s.invRepo.Insert(inv); err != nil {
		return nil, fmt.Errorf("failed to create invite: %v", err)
	}

	// TODO - send out email with invitation link

	return inv, nil
}

// ValidateInvite loads the Invite for the provided email/token combination and
// returns the Invite's unique identifier if it is active.
// Returns an error if the Invite does not exist for the provided email/token
// combination, or the Invite has expired.
//
// email is the email address to look up the Invite for.
//
// token is the token string that should match with the email on the Invite.
func (s *InviteService) ValidateInvite(email string, token string) (types.InviteID, error) {
	inv, err := s.invRepo.Get(email, token)
	if err != nil {
		return "", err
	}

	if inv == nil {
		return "", &types.InviteNotFoundError{Email: email, Token: token}
	}

	if inv.Expiration.Before(time.Now().UTC()) {
		return "", &types.InviteExpiredError{Email: email, Token: token}
	}

	if inv.Accepted {
		return "", &types.InviteAlreadyAcceptedError{Email: email, Token: token}
	}

	return inv.ID, nil
}

// AcceptInvite marks the Invite with the provided ID as accepted.
//
// id is the unique identifier of the Invite being accepted.
func (s *InviteService) AcceptInvite(id types.InviteID) error {
	return s.invRepo.Accept(id)
}
