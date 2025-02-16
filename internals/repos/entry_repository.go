package repos

import (
	"fmt"

	"github.com/mhs294/mulhall/internals/db"
	"github.com/mhs294/mulhall/internals/types"
	"go.mongodb.org/mongo-driver/bson"
)

// EntryRepository manages Entry records in the database.
type EntryRepository struct {
	mdb      *db.MongoDB
	dbName   string
	collName string
}

// NewEntryRepository creates a new EntryRepository instance and returns a pointer to it.
//
// mdb is the MongoDB instance used by the EntryRepository.
func NewEntryRepository(mdb *db.MongoDB) *EntryRepository {
	return &EntryRepository{mdb: mdb, dbName: "mulhall", collName: "entries"}
}

// TestConnection tests the connection to the MongoDB instance and returns any error that occurs.
func (r *EntryRepository) TestConnection() error {
	return r.mdb.TestConnection(r.dbName)
}

// Insert inserts the provided Entry into the database.
//
// e is the Entry to insert into the database.
func (r *EntryRepository) Insert(e *types.Entry) error {
	if err := r.mdb.InsertOne(r.dbName, r.collName, e); err != nil {
		return fmt.Errorf("failed to insert entry: %v", err)
	}

	return nil
}

// GetBySchedule gets all Entries for the provided Schedule.
//
// id is the unique identifier of the Schedule for which corresponding Entries should be loaded.
func (r *EntryRepository) GetBySchedule(id types.ScheduleID) ([]types.Entry, error) {
	// Define the query
	query := bson.M{"schedule": id}

	// Load Entries from the database
	var entries []types.Entry
	if err := r.mdb.GetAll(r.dbName, r.collName, query, &entries); err != nil {
		return nil, fmt.Errorf("failed to look up entries (schedule=%v): %v", id, err)
	}

	return entries, nil
}

// GetByContestant gets all Entries for the provided Contestant.
//
// id is the unique identifier of the Contestant for which corresponding Entries should be loaded.
func (r *EntryRepository) GetByContestant(id types.ContestantID) ([]types.Entry, error) {
	// Define the query
	query := bson.M{"contestant": id}

	// Load Entries from the database
	var entries []types.Entry
	if err := r.mdb.GetAll(r.dbName, r.collName, query, &entries); err != nil {
		return nil, fmt.Errorf("failed to look up entries (contestant=%v): %v", id, err)
	}

	return entries, nil
}

// Update updates a Entry in the database using the information in the provided model.
//
// e is the model to use to update the Entry in the database. The models' EntryID is used to
// determine which document in the database should be replaced with the updated version.
func (r *EntryRepository) Update(e *types.Entry) error {
	// Define the filter query
	filter := bson.M{"id": e.ID}

	// Perform the update
	if err := r.mdb.ReplaceOne(r.dbName, r.collName, filter, e); err != nil {
		return fmt.Errorf("failed to update entry (%v): %v", e, err)
	}

	return nil
}
