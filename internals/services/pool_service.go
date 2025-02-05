package services

import (
	"fmt"

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
		Contestants: make(map[types.ContestantID]struct{}, 0),
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
func (s *PoolService) AddContestant(poolID types.PoolID, conID types.ContestantID) error {
	// Load the Pool from the database
	p, err := s.repo.GetByID(poolID)
	if err != nil {
		return err
	}

	// If the Contestant is already in the Pool, do nothing and return
	if _, exists := p.Contestants[conID]; exists {
		return nil
	}

	// Add the Contestant to the Pool
	p.Contestants[conID] = struct{}{}
	if err = s.repo.Update(p); err != nil {
		return fmt.Errorf("failed to add contestant to pool (pool=%s, contestant=%s): %v", poolID, conID, err)
	}

	return nil
}

// RemoveContestant removes the specified Contestant from the specified Pool.
//
// poolID is the unique identifier of the Pool to update.
//
// conID is the unique identifier of the Contestant to remove from the Pool.
func (s *PoolService) RemoveContestant(poolID types.PoolID, conID types.ContestantID) error {
	// Load the Pool from the database
	p, err := s.repo.GetByID(poolID)
	if err != nil {
		return err
	}

	// If the Contestant is not in the Pool, do nothing and return
	if _, exists := p.Contestants[conID]; !exists {
		return nil
	}

	// Remove the Contestant from the Pool
	delete(p.Contestants, conID)
	if err = s.repo.Update(p); err != nil {
		return fmt.Errorf("failed to remove contestant to pool (pool=%s, contestant=%s): %v", poolID, conID, err)
	}

	return nil
}

// Complete marks the specified Pool as complete (i.e. - its contest has concluded)
//
// id is the unique identifier of the Pool to mark as complete.
func (s *PoolService) Complete(id types.PoolID) error {
	// Load the Pool from the database
	p, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// If the Pool is already marked as complete, do nothing and return
	if p.Complete {
		return nil
	}

	// Mark the Pool as complete
	p.Complete = true
	if err = s.repo.Update(p); err != nil {
		return fmt.Errorf("failed to mark pool as complete (id=%s): %v", id, err)
	}

	return nil
}

// Deactivate deactivates the specified Pool (soft-delete).
//
// id is the unique identifier of the Pool to deactivate.
func (s *PoolService) Deactivate(id types.PoolID) error {
	return s.repo.Deactivate(id)
}
