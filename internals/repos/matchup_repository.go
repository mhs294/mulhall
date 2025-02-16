package repos

import (
	"fmt"

	"github.com/mhs294/mulhall/internals/db"
	"github.com/mhs294/mulhall/internals/types"
	"go.mongodb.org/mongo-driver/bson"
)

// MatchupRepository manages Matchup records in the database.
type MatchupRepository struct {
	mdb      *db.MongoDB
	dbName   string
	collName string
}

// NewMatchupRepository creates a new MatchupRepository instance and returns a pointer to it.
//
// mdb is the MongoDB instance used by the MatchupRepository.
func NewMatchupRepository(mdb *db.MongoDB) *MatchupRepository {
	return &MatchupRepository{mdb: mdb, dbName: "mulhall", collName: "matchups"}
}

// TestConnection tests the connection to the MongoDB instance and returns any error that occurs.
func (r *MatchupRepository) TestConnection() error {
	return r.mdb.TestConnection(r.dbName)
}

// Insert inserts the provided Matchup into the database.
//
// m is the Matchup to insert into the database.
func (r *MatchupRepository) Insert(m *types.Matchup) error {
	if err := r.mdb.InsertOne(r.dbName, r.collName, m); err != nil {
		return fmt.Errorf("failed to insert matchup: %v", err)
	}

	return nil
}

// GetByIDs returns all Matchups in the database for the specified IDs.
//
// ids is the slice of unique identifiers of the Matchups to load.
func (r *MatchupRepository) GetByIDs(ids []types.MatchupID) ([]types.Matchup, error) {
	// Define the query
	query := bson.M{"id": bson.M{"$in": ids}}

	// Load the Matchups from the database
	var matchups []types.Matchup
	if err := r.mdb.GetAll(r.dbName, r.collName, query, &matchups); err != nil {
		return nil, fmt.Errorf("failed to look up matchups (ids=%v): %v", ids, err)
	}

	return matchups, nil
}

// Update updates a Matchup in the database using the information in the provided model.
//
// e is the model to use to update the Matchup in the database. The models' MatchupID is used to
// determine which document in the database should be replaced with the updated version.
func (r *MatchupRepository) Update(m *types.Matchup) error {
	// Define the filter query
	filter := bson.M{"id": m.ID}

	// Perform the update
	if err := r.mdb.ReplaceOne(r.dbName, r.collName, filter, m); err != nil {
		return fmt.Errorf("failed to update matchup (%v): %v", m, err)
	}

	return nil
}
