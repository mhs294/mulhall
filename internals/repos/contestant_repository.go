package repos

import (
	"fmt"

	"github.com/mhs294/mulhall/internals/db"
	"github.com/mhs294/mulhall/internals/types"
	"go.mongodb.org/mongo-driver/bson"
)

// ContestantRepository manages Contestant records in the database.
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

// Insert inserts the provided Contestant into the database.
//
// c is the Contestant to insert into the database.
func (r *ContestantRepository) Insert(c *types.Contestant) error {
	if err := r.mdb.InsertOne(r.dbName, r.collName, c); err != nil {
		return fmt.Errorf("failed to insert contestant: %v", err)
	}

	return nil
}

// GetByIDs returns all active Contestants in the database for the specified IDs.
//
// ids is the slice of unique identifiers load Contestants for.
func (r *ContestantRepository) GetByIDs(ids []types.ContestantID) ([]types.Contestant, error) {
	// Define the query
	query := bson.M{
		"id":     bson.M{"$in": ids},
		"active": true,
	}

	// Load the Contestants from the database
	var cons []types.Contestant
	if err := r.mdb.GetAll(r.dbName, r.collName, query, &cons); err != nil {
		return nil, fmt.Errorf("failed to look up contestants (ids=%v): %v", ids, err)
	}

	return cons, nil
}

// GetByAuthorizedUser returns all active Contestants in the database for which the specified User is authorized
// (or an empty slice if the User is not authorized for any Contestants).
//
// userID is the unique identifier of the authorized User to load Contestants for.
func (r *ContestantRepository) GetByAuthorizedUser(userID types.UserID) ([]types.Contestant, error) {
	// Define the query
	query := bson.M{
		fmt.Sprintf("authorizedUsers.%s", userID): bson.M{"$exists": "true"},
		"active": true,
	}

	// Load the Contestants from the database
	var cons []types.Contestant
	if err := r.mdb.GetAll(r.dbName, r.collName, query, &cons); err != nil {
		return nil, fmt.Errorf("failed to look up contestants (user id=%s): %v", userID, err)
	}

	return cons, nil
}

// Update updates a Contestant in the database using the information in the provided model.
//
// c is the model to use to update the Contestant in the database. The models' ContestantID is used to
// determine which document in the database should be replaced with the updated version.
func (r *ContestantRepository) Update(c *types.Contestant) error {
	// Define the filter query
	filter := bson.M{
		"id":     c.ID,
		"active": true,
	}

	// Perform the update
	if err := r.mdb.ReplaceOne(r.dbName, r.collName, filter, c); err != nil {
		return fmt.Errorf("failed to update contestant (%v): %v", c, err)
	}

	return nil
}

// Deactivate sets the Contestant with the provided ID to be inactive.
//
// id the unique identifier of the Contestant to deactivate.
func (r *ContestantRepository) Deactivate(id types.ContestantID) error {
	// Define the filter query and update operation
	filter := bson.M{"id": id}
	update := bson.M{
		"$set": bson.M{
			"active": false,
		},
	}

	// Perform the update
	if err := r.mdb.UpdateOne(r.dbName, r.collName, filter, update); err != nil {
		return fmt.Errorf("failed to deactivate contestant (id=%s): %v", id, err)
	}

	return nil
}
