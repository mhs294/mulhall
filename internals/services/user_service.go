package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mhs294/mulhall/internals/env"
	"github.com/mhs294/mulhall/internals/repos"
	"github.com/mhs294/mulhall/internals/types"
	"github.com/mhs294/mulhall/internals/utils"
	"golang.org/x/crypto/bcrypt"
)

// UserService represents a service for interacting with Users and their accounts on the site.
type UserService struct {
	invService *InviteService
	userRepo   *repos.UserRepository
	sessRepo   *repos.SessionRepository
}

// NewUserService creates a new instance of an UserService and returns a pointer to it.
//
// s is the InviteService that will be used to manage Invitations during User creation workflows.
//
// r is the UserRepository that will be used at runtime by the UserService.
func NewUserService(s *InviteService, ur *repos.UserRepository, sr *repos.SessionRepository) *UserService {
	return &UserService{
		invService: s,
		userRepo:   ur,
		sessRepo:   sr,
	}
}

// Register handles the creation of a new User from an accepted Invite.
//
// req is the RegisterUserRequest containing the information required to create the User and accept the Invite.
func (s *UserService) Register(req *types.RegisterUserRequest) (*types.User, error) {
	// Validate the request
	if req.Password != req.Confirm {
		return nil, &types.PasswordMismatchError{}
	}
	invId, err := s.invService.ValidateInvite(req.Email, req.Token)
	if err != nil {
		return nil, err
	}

	// Create the User
	u, err := s.createUser(req)
	if err != nil {
		return nil, err
	}

	// Mark the Invite as accepted
	if err = s.invService.AcceptInvite(invId); err != nil {
		// TODO - figure out how rollbacks should be handled
		return nil, err
	}

	return u, nil
}

// Login authenticates a User from the provided email and password and returns a new
// Session for that User if authentication succeeds. If login fails, an error is returned.
//
// email is the User's email address
//
// pwd is the raw password submitted by the User.
func (s *UserService) Login(email string, pwd string) (*types.Session, error) {
	// Look up the User
	u, err := s.userRepo.GetUser(email)
	if err != nil {
		var notFound *types.UserNotFoundError
		if errors.As(err, &notFound) {
			return nil, notFound
		}

		return nil, fmt.Errorf("failed to login user: %v", err)
	}

	// Compare the submitted password hash against the User's password hash
	hash, err := hashPassword(pwd, u.Salt)
	if err != nil {
		return nil, fmt.Errorf("failed to login user: %v", err)
	}
	if hash != u.Hash {
		return nil, &types.PasswordIncorrectError{}
	}

	// Authentication successful, create new Session for the User
	sess := &types.Session{
		ID:         types.SessionID(uuid.New().String()),
		UserID:     u.ID,
		Expiration: time.Now().UTC().Add(env.SessionExpiration),
	}
	if err = s.sessRepo.InsertSession(sess); err != nil {
		return nil, fmt.Errorf("failed to create new session for user: %v", err)
	}

	return sess, nil
}

func (s *UserService) createUser(req *types.RegisterUserRequest) (*types.User, error) {
	// Create the User with a new randomly generated salt and hash the provided password with it
	salt := utils.CreateToken(16)
	hash, err := hashPassword(req.Password, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}
	u := &types.User{
		ID:            types.UserID(uuid.New().String()),
		Email:         req.Email,
		Salt:          salt,
		Hash:          hash,
		Administrator: false,
		Active:        true,
	}

	// Insert the User into the database
	if err := s.userRepo.InsertUser(u); err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return u, nil
}

func hashPassword(pwd string, salt string) (string, error) {
	bytes := []byte(pwd + salt)
	hash, err := bcrypt.GenerateFromPassword(bytes, bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
