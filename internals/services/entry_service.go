package services

import "github.com/mhs294/mulhall/internals/repos"

// EntryService represents a service for managing Contestants' Entries in a Pool.
type EntryService struct {
	repo *repos.EntryRepository
}

// NewEntryService creates a new instance of a EntryService and returns a pointer to it.
//
// r is the EntryRepository used to manage Entry records in the database.
func NewEntryService(r *repos.EntryRepository) *EntryService {
	return &EntryService{repo: r}
}

// TODO - Create

// TODO - GetByID

// TODO - SetSelectedPick

// TODO - ClearSelectedPick

// TODO - AddSuggestedPick

// TODO - RemoveSuggestedPick

// TODO - Deactivate
