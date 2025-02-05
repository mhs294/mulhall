package services

import (
	"github.com/google/uuid"
	"github.com/mhs294/mulhall/internals/repos"
	"github.com/mhs294/mulhall/internals/types"
)

type PoolService struct {
	repo *repos.PoolRepository
}

// NewPoolService creates a new instance of a PoolService and returns a pointer to it.
//
// r is the PoolRepository that will be used to manage Pool records within the database.
func NewPoolService(r *repos.PoolRepository) *PoolService {
	return &PoolService{repo: r}
}

// Create creates a new Pool from the provided information.
// Returns an updated version of the Pool model containing its ID after creation.
//
// req is the CreatePoolRequest containing the information required to create the Pool.
func (s *PoolService) Create(req *types.CreatePoolRequest) (*types.Pool, error) {
	p := &types.Pool{
		ID:          types.PoolID(uuid.NewString()),
		Name:        req.Name,
		Contestants: make([]types.ContestantID, 0),
		Active:      true,
		Complete:    false,
	}

	if err := s.repo.Insert(p); err != nil {
		return nil, err
	}

	return p, nil
}

// GetAll returns all active Pools.
func (s *PoolService) GetAll() ([]types.Pool, error) {
	return s.repo.GetAll()
}

// GetByID gets the Pool for the provided ID.
//
// id is the unique identifier of the Pool to look up.
func (s *PoolService) GetByID(id types.PoolID) (*types.Pool, error) {
	return s.repo.GetByID(id)
}

// AddContestant adds the specified Contestant to the specified Pool.
//
// poolID is the unique identifier of the Pool to update.
//
// conID is the unique identifier of the Contestant to add to the Pool.
func (s *PoolService) AddContestant(poolId types.PoolID, conID types.ContestantID) error {
	return s.repo.AddContestant(poolId, conID)
}

// RemoveContestant removes the specified Contestant from the specified Pool.
//
// poolID is the unique identifier of the Pool to update.
//
// conID is the unique identifier of the Contestant to remove from the Pool.
func (s *PoolService) RemoveContestant(poolId types.PoolID, conID types.ContestantID) error {
	return s.repo.RemoveContestant(poolId, conID)
}

// Complete marks the specified Pool as complete (i.e. - its contest has concluded)
//
// id is the unique identifier of the Pool to mark as complete.
func (s *PoolService) Complete(id types.PoolID) error {
	return s.repo.Complete(id)
}

// TODO - deactivate
