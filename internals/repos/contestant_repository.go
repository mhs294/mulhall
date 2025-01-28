package repos

import (
	"fmt"

	"github.com/mhs294/mulhall/internals/db"
	"github.com/mhs294/mulhall/internals/types"
	"go.mongodb.org/mongo-driver/bson"
)

// ContestantRepository is a mechanism for managing Contestants participating in Pools on the site.
type ContestantRepository struct {
	mdb      *db.MongoDB
	dbName   string
	collName string
}

// NewContestantRepository creates a new ContestantRepository instance and returns a pointer to it.
//
// mdb is the MongoDB instance used by the ContestantRepository.
func NewContestantRepository(mdb *db.MongoDB) *ContestantRepository {
	return &ContestantRepository{mdb: mdb, dbName: "mulhall", collName: "contestants"}
}

// TestConnection tests the connection to the MongoDB instance and returns any error that occurs.
func (r *ContestantRepository) TestConnection() error {
	return r.mdb.TestConnection(r.dbName)
}

// InsertContestant inserts the provided Contestant into the database.
//
// c is the Contestant to insert into the database.
func (r *ContestantRepository) InsertContestant(c *types.Contestant) error {
	if err := r.mdb.InsertOne(r.dbName, r.collName, c); err != nil {
		return fmt.Errorf("failed to insert contestant: %v", err)
	}

	return nil
}

// GetContestant returns the Contestant for the provided ID (or nil if no such Contestant exists).
//
// id is the unique identifier of the Contestant to look up.
func (r *ContestantRepository) GetContestant(id types.ContestantID) (*types.Contestant, error) {
	var sess types.Contestant
	if err := r.mdb.GetOne(r.dbName, r.collName, bson.D{{Key: "id", Value: id}}, &sess); err != nil {
		return nil, fmt.Errorf("failed to look up contestant: %v", err)
	}

	return &sess, nil
}

// UpdateContestant updates a Contestant in the database using the information in the provided model.
//
// c is the model to use to update the Contestant in the database. The models' ContestantID is used to
// determine which document in the database should be replaced with the updated version.
func (r *ContestantRepository) UpdateContestant(c *types.Contestant) error {
	// Define the filter query
	filter := bson.M{"id": c.ID}

	// Perform the update
	if err := r.mdb.ReplaceOne(r.dbName, r.collName, filter, c); err != nil {
		return fmt.Errorf("failed to update contestant (%v): %v", c, err)
	}

	return nil
}
