package services

import (
	"github.com/google/uuid"
	"github.com/mhs294/mulhall/internals/repos"
	"github.com/mhs294/mulhall/internals/types"
	"github.com/mhs294/mulhall/internals/types/roles"
	"github.com/mhs294/mulhall/internals/types/status"
)

// ContestantService represents a service for managing Contestants and their authorized Users.
type ContestantService struct {
	repo        *repos.ContestantRepository
	poolService *PoolService
}

// NewContestantService creates a new instance of a ContestantService and returns a pointer to it.
//
// r is the ContestantRepository that will be used to manage Contestant records in the database.
//
// ps is the PoolService that will be used to load Contestant information for specific Pools.
func NewContestantService(r *repos.ContestantRepository, ps *PoolService) *ContestantService {
	return &ContestantService{repo: r, poolService: ps}
}

// GetByPool returns all active Contestants for the specified Pool.
//
// poolID is the unique identifier of the Pool to load Contestants for.
func (s *ContestantService) GetByPool(poolID types.PoolID) ([]types.Contestant, error) {
	// Load the Pool
	p, err := s.poolService.GetByID(poolID)
	if err != nil {
		return nil, err
	}

	// Load Contestants from database using Contestant IDs from Pool
	i := 0
	conIDs := make([]types.ContestantID, len(p.Contestants))
	for id := range p.Contestants {
		conIDs[i] = id
		i++
	}
	c, err := s.repo.GetByIDs(conIDs)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// GetByAuthorizedUser returns all active Contestants for which the specified User is authorized.
//
// userID is the unique identifier of the authorized User to load Contestants for.
func (s *ContestantService) GetByAuthorizedUser(userID types.UserID) ([]types.Contestant, error) {
	return s.repo.GetByAuthorizedUser(userID)
}

// Create creates a new Contestant from the provided information.
// Returns an updated version of the Contestant model containing its ID after creation.
//
// req is the CreateContestantRequest containing the information required to create the Contestant.
func (s *ContestantService) Create(req *types.CreateContestantRequest) (*types.Contestant, error) {
	c := &types.Contestant{
		ID:              types.ContestantID(uuid.NewString()),
		Name:            req.Name,
		AuthorizedUsers: req.AuthorizedUsers,
		Active:          true,
		Status:          status.ACTIVE,
	}

	if err := s.repo.Insert(c); err != nil {
		return nil, err
	}

	return c, nil
}

// SetAuthorizedUser sets the specified User to be authorized for the Contestant with the specified Role.
//
// conID is the unique identifier of the Contestant to update.
//
// userID is the unique identifier of the User to authorize for the Contestant.
//
// role is the Role for the User that will dictate its level of access to manage the Contestant.
func (s *ContestantService) SetAuthorizedUser(conID types.ContestantID, userID types.UserID, role roles.Role) error {
	// Load the Contestant from the database
	c, err := s.repo.GetByID(conID)
	if err != nil {
		return err
	}

	// Update the authorized User's Role for the Contestant
	c.AuthorizedUsers[userID] = role
	if err = s.repo.Update(c); err != nil {
		return err
	}

	return nil
}

// RemoveAuthorizedUser removed the specified User from the list of authorized Users for the Contestant.
//
// conID is the unique identifier of the Contestant to update.
//
// userID is the unique identifier of the authorized User to remove from the Contestant.
func (s *ContestantService) RemoveAuthorizedUser(conID types.ContestantID, userID types.UserID) error {
	// Load the Contestant from the database
	c, err := s.repo.GetByID(conID)
	if err != nil {
		return err
	}

	// Remove the authorized User from the Contestant
	delete(c.AuthorizedUsers, userID)
	if err = s.repo.Update(c); err != nil {
		return err
	}

	return nil
}

// SetStatus updates the Status of the specified Contestant.
//
// id is the unique identifier of the Contestant to update.
//
// status is the new Status that will be applied to the Contestant.
func (s *ContestantService) SetStatus(id types.ContestantID, status status.Status) error {
	// Load the Contestant from the database
	c, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// Update the Contestant's Status
	c.Status = status
	if err = s.repo.Update(c); err != nil {
		return err
	}

	return nil
}

// Deactivate deactivates the specified Contestant (soft-delete).
//
// id is the unique identifier of the Contestant to deactivate.
func (s *ContestantService) Deactivate(id types.ContestantID) error {
	return s.repo.Deactivate(id)
}
