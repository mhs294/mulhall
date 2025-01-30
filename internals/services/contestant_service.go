package services

import (
	"github.com/mhs294/mulhall/internals/repos"
	"github.com/mhs294/mulhall/internals/types"
)

// ContestantService represents a service for managing Contestants and their authorized Users.
type ContestantService struct {
	repo *repos.ContestantRepository
}

// NewContestantService creates a new instance of a ContestantService and returns a pointer to it.
//
// r is the ContestantRepository that will be used at runtime by the ContestantService.
func NewContestantService(r *repos.ContestantRepository) *ContestantService {
	return &ContestantService{repo: r}
}

// GetByPool returns all active Contestants for the specified Pool.
//
// poolID is the unique identifier of the Pool to load Contestants for.
func (s *ContestantService) GetByPool(poolID *types.PoolID) ([]types.Contestant, error) {
	// TODO - start here
}

// TODO - stub remaining methods
