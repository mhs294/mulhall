package repos

import (
	"fmt"
	"sort"

	"github.com/mhs294/mulhall/internals/db"
	"github.com/mhs294/mulhall/internals/types"
	"go.mongodb.org/mongo-driver/bson"
)

// TeamRepository manages Team records in the database.
// Since NFL teams rarely change, the Teams are loaded on initialization
// and then cached for the lifetime of the application.
type TeamRepository struct {
	mdb      *db.MongoDB
	dbName   string
	collName string
	teams    map[types.TeamID]types.Team
}

// NewTeamRepository creates a new TeamRepository instance and returns a pointer to it.
// connStr is the connection string for the MongoDB instance (e.g. - "mongodb+srv://{user}:{pass}@myinstance.mongodb.net/").
func NewTeamRepository(mdb *db.MongoDB) *TeamRepository {
	return &TeamRepository{mdb: mdb, dbName: "mulhall", collName: "teams", teams: nil}
}

// TestConnection tests the connection to the MongoDB instance and returns any error that occurs.
func (r *TeamRepository) TestConnection() error {
	return r.mdb.TestConnection(r.dbName)
}

// GetAll returns a slice of all available Teams, sorted by their location shorthand.
func (r *TeamRepository) GetAll() ([]types.Team, error) {
	if r.teams == nil {
		err := r.loadTeams()
		if err != nil {
			return nil, err
		}
	}

	teams := make([]types.Team, 0, len(r.teams))
	for _, t := range r.teams {
		teams = append(teams, t)
	}
	sort.Slice(teams, func(i, j int) bool {
		return teams[i].Shorthand < teams[j].Shorthand
	})

	return teams, nil
}

// GetByID returns the Team keyed by the specified ID (or an empty Team if no such Team exists).
//
// id is the unique identifier of the Team to look up.
func (r *TeamRepository) GetByID(id types.TeamID) (types.Team, error) {
	if r.teams == nil {
		err := r.loadTeams()
		if err != nil {
			return types.Team{}, err
		}
	}

	return r.teams[id], nil
}

func (r *TeamRepository) loadTeams() error {
	var teams []types.Team
	if err := r.mdb.GetAll(r.dbName, r.collName, bson.M{}, &teams); err != nil {
		return fmt.Errorf("failed to load teams from database: %v", err)
	}

	r.teams = make(map[types.TeamID]types.Team, len(teams))
	for _, t := range teams {
		r.teams[t.ID] = t
	}

	return nil
}
