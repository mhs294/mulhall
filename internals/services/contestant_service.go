package services

import (
	"github.com/google/uuid"
	"github.com/mhs294/mulhall/internals/repos"
	"github.com/mhs294/mulhall/internals/types"
	"github.com/mhs294/mulhall/internals/types/status"
)

// ContestantService represents a service for managing Contestants and their authorized Users.
type ContestantService struct {
	conRepo     *repos.ContestantRepository
	poolService *PoolService
}

// NewContestantService creates a new instance of a ContestantService and returns a pointer to it.
//
// cr is the ContestantRepository that will be used to manage Contestant records in the database.
//
// ps is the PoolService that will be used to load Contestant information for specific Pools.
func NewContestantService(cr *repos.ContestantRepository, ps *PoolService) *ContestantService {
	return &ContestantService{conRepo: cr, poolService: ps}
}

// GetByPool returns all active Contestants for the specified Pool.
//
// poolID is the unique identifier of the Pool to load Contestants for.
func (s *ContestantService) GetByPool(poolID types.PoolID) ([]types.Contestant, error) {
	p, err := s.poolService.GetByID(poolID)
	if err != nil {
		return nil, err
	}

	c, err := s.conRepo.GetByIDs(p.Contestants)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// GetByAuthorizedUser returns all active Contestants for which the specified User is authorized.
//
// userID is the unique identifier of the authorized User to load Contestants for.
func (s *ContestantService) GetByAuthorizedUser(userID types.UserID) ([]types.Contestant, error) {
	return s.conRepo.GetByAuthorizedUser(userID)
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

	if err := s.conRepo.Insert(c); err != nil {
		return nil, err
	}

	return c, nil
}

// TODO - add authorized user

// TODO - remove authorized user

// TODO - deactivate
