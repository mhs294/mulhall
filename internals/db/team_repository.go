package db

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/mhs294/mulhall/internals/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const db = "mulhall"
const teamsColl = "teams"

// TeamRepository is a mechanism for loading Team information.
// Since NFL teams rarely change, the Teams are loaded on initialization
// and then cached for the lifetime of the application.
type TeamRepository struct {
	teams map[string]types.Team
}

// NewTeamRepository creates a new TeamRepository instance and returns a pointer to it.
func NewTeamRepository(connStr string) (*TeamRepository, error) {
	// Load Teams from database
	teams, err := loadTeams(connStr)
	if err != nil {
		return nil, err
	}

	// Initialize TeamRepository internal map
	tr := TeamRepository{teams: make(map[string]types.Team, len(teams))}
	for _, t := range teams {
		tr.teams[t.ID] = t
	}

	return &tr, nil
}

// loadTeams loads all Teams from the database and returns them in a slice.
// If an error occurs during load, that error will be returned alongside a nil slice.
func loadTeams(connStr string) ([]types.Team, error) {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(connStr).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	// Setup deferred connection closure for when function completes
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// Load all of the documents from the teams collection
	coll := client.Database(db).Collection(teamsColl)
	cursor, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to query %s.%s: %v", db, teamsColl, err)
	}

	// Unpack the cursor contents into a slice
	var teams []types.Team
	if err = cursor.All(ctx, &teams); err != nil {
		return nil, fmt.Errorf("failed to parse the query results: %v", err)
	}

	// Initialization successful
	return teams, nil
}

// GetAllTeams returns a slice of all available Teams, sorted by their location shorthand.
func (tr *TeamRepository) GetAllTeams() []types.Team {
	teams := make([]types.Team, 0, len(tr.teams))
	for _, t := range tr.teams {
		teams = append(teams, t)
	}
	sort.Slice(teams, func(i, j int) bool {
		return teams[i].Shorthand < teams[j].Shorthand
	})

	return teams
}

// GetTeam returns the Team keyed by the specified ID (or nil if no such Team exists).
func (tr *TeamRepository) GetTeam(id string) types.Team {
	return tr.teams[id]
}
