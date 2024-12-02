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
	InviteRepository *db.InviteRepository
}

// NewInviteService creates a new instance of an InviteService and returns a pointer to it.
// ir is the pointer to the InviteRepository that will be used at runtime by the InviteService.
func NewInviteService(ir *db.InviteRepository) *InviteService {
	return &InviteService{InviteRepository: ir}
}

// CreateInvite creates a new Invite from the provided information and returns a pointer to it.
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
	if err := s.InviteRepository.InsertInvite(invite); err != nil {
		return nil, fmt.Errorf("failed to create invite: %v", err)
	}

	return invite, nil
}

// TODO - start here
