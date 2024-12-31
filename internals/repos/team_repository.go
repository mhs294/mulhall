package repos

import (
	"fmt"
	"sort"

	"github.com/mhs294/mulhall/internals/db"
	"github.com/mhs294/mulhall/internals/types"
	"go.mongodb.org/mongo-driver/bson"
)

// TeamRepository is a mechanism for loading Team information.
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
func (tr *TeamRepository) TestConnection() error {
	return tr.mdb.TestConnection(tr.dbName)
}

// GetAllTeams returns a slice of all available Teams, sorted by their location shorthand.
func (tr *TeamRepository) GetAllTeams() ([]types.Team, error) {
	if tr.teams == nil {
		err := tr.loadTeams()
		if err != nil {
			return nil, err
		}
	}

	teams := make([]types.Team, 0, len(tr.teams))
	for _, t := range tr.teams {
		teams = append(teams, t)
	}
	sort.Slice(teams, func(i, j int) bool {
		return teams[i].Shorthand < teams[j].Shorthand
	})

	return teams, nil
}

// GetTeam returns the Team keyed by the specified ID (or an empty Team if no such Team exists).
func (tr *TeamRepository) GetTeam(id types.TeamID) (types.Team, error) {
	if tr.teams == nil {
		err := tr.loadTeams()
		if err != nil {
			return types.Team{}, err
		}
	}

	return tr.teams[id], nil
}

func (tr *TeamRepository) loadTeams() error {
	var teams []types.Team
	if err := tr.mdb.GetAll(tr.dbName, tr.collName, bson.D{}, &teams); err != nil {
		return fmt.Errorf("failed to load teams from database: %v", err)
	}

	tr.teams = make(map[types.TeamID]types.Team, len(teams))
	for _, t := range teams {
		tr.teams[t.ID] = t
	}

	return nil
}
