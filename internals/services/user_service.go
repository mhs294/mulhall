package services

import "github.com/mhs294/mulhall/internals/repos"

// UserService represents a service for interacting with Users and their accounts on the site.
type UserService struct {
	userRepository *repos.UserRepository
}

// NewUserService creates a new instance of an UserService and returns a pointer to it.
//
// r is the UserRepository that will be used at runtime by the UserService.
func NewUserService(r *repos.UserRepository) *UserService {
	return &UserService{userRepository: r}
}

// TODO - start here
