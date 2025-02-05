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

// NewUserService creates a new instance of a UserService and returns a pointer to it.
//
// s is the InviteService that will be used to manage Invitations during User creation workflows.
//
// ur is the UserRepository that will be used to manage User records in the database.
//
// sr is the SessionRepository that will be used to manage Session records in the database.
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
	invId, err := s.invService.Validate(req.Email, req.Token)
	if err != nil {
		return nil, err
	}

	// Create the User
	u, err := s.createUser(req)
	if err != nil {
		return nil, err
	}

	// Mark the Invite as accepted
	if err = s.invService.Accept(invId); err != nil {
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
	u, err := s.userRepo.GetByEmail(email)
	if err != nil {
		var notFound *types.UserNotFoundError
		if errors.As(err, &notFound) {
			return nil, notFound
		}

		return nil, fmt.Errorf("failed to login user: %v", err)
	}

	// Compare the submitted password against the User's password
	saltedPwd := []byte(pwd + u.Salt)
	err = bcrypt.CompareHashAndPassword([]byte(u.Hash), saltedPwd)
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, &types.PasswordIncorrectError{}
	} else if err != nil {
		return nil, fmt.Errorf("failed to login user: %v", err)
	}

	// Authentication successful, create new Session for the User
	sess := &types.Session{
		ID:         types.SessionID(uuid.NewString()),
		User:       u.ID,
		Expiration: time.Now().UTC().Add(env.SessionExpiration),
	}
	if err = s.sessRepo.Insert(sess); err != nil {
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
		ID:            types.UserID(uuid.NewString()),
		Email:         req.Email,
		Salt:          salt,
		Hash:          hash,
		Administrator: false,
		Active:        true,
	}

	// Insert the User into the database
	if err := s.userRepo.Insert(u); err != nil {
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
