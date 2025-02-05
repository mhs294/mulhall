package repos

import (
	"fmt"

	"github.com/mhs294/mulhall/internals/db"
	"github.com/mhs294/mulhall/internals/types"
	"go.mongodb.org/mongo-driver/bson"
)

// PoolRepository manages Pool records in the database.
type PoolRepository struct {
	mdb      *db.MongoDB
	dbName   string
	collName string
}

// NewPoolRepository creates a new PoolRepository instance and returns a pointer to it.
//
// mdb is the MongoDB instance used by the PoolRepository.
func NewPoolRepository(mdb *db.MongoDB) *PoolRepository {
	return &PoolRepository{mdb: mdb, dbName: "mulhall", collName: "contestants"}
}

// TestConnection tests the connection to the MongoDB instance and returns any error that occurs.
func (r *PoolRepository) TestConnection() error {
	return r.mdb.TestConnection(r.dbName)
}

// Insert inserts the provided Pool into the database.
//
// p is the Pool to insert into the database.
func (r *PoolRepository) Insert(p *types.Pool) error {
	if err := r.mdb.InsertOne(r.dbName, r.collName, p); err != nil {
		return fmt.Errorf("failed to insert pool: %v", err)
	}

	return nil
}

// GetAll loads all active Pools from the database.
func (r *PoolRepository) GetAll() ([]types.Pool, error) {
	var pools []types.Pool
	if err := r.mdb.GetAll(r.dbName, r.collName, bson.M{}, &pools); err != nil {
		return nil, fmt.Errorf("failed to load pools from database: %v", err)
	}

	return pools, nil
}

// GetByID gets the Pool for the provided ID.
// Returns UserNotFoundError if no such User exists or has been deactivated.
//
// id is the unique identifier of the Pool to look up.
func (r *PoolRepository) GetByID(id types.PoolID) (*types.Pool, error) {
	// Define the query
	query := bson.M{"id": id}

	// Load Pool from the database
	var p types.Pool
	if err := r.mdb.GetOne(r.dbName, r.collName, query, &p); err != nil {
		return nil, fmt.Errorf("failed to look up pool: %v", err)
	}

	// Verify that the Pool exists and is active
	if p.ID != id || !p.Active {
		return nil, &types.PoolNotFoundError{ID: id}
	}

	return &p, nil
}

// Update updates a Pool in the database using the information in the provided model.
//
// p is the model to use to update the Pool in the database. The models' PoolID is used to
// determine which document in the database should be replaced with the updated version.
func (r *PoolRepository) Update(p *types.Pool) error {
	// Define the filter query
	filter := bson.M{
		"id":     p.ID,
		"active": true,
	}

	// Perform the update
	if err := r.mdb.ReplaceOne(r.dbName, r.collName, filter, p); err != nil {
		return fmt.Errorf("failed to update pool (%v): %v", p, err)
	}

	return nil
}

// Deactivate sets the Pool with the provided ID to be inactive.
//
// id the unique identifier of the Pool to deactivate.
func (r *PoolRepository) Deactivate(id types.PoolID) error {
	// Define the filter query and update operation
	filter := bson.M{"id": id}
	update := bson.M{
		"$set": bson.M{
			"active": false,
		},
	}

	// Perform the update
	if err := r.mdb.UpdateOne(r.dbName, r.collName, filter, update); err != nil {
		return fmt.Errorf("failed to deactivate pool (id=%s): %v", id, err)
	}

	return nil
}
