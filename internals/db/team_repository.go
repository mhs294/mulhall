package db

import (
	"context"
	"fmt"
	"log"
	"sort"

	"github.com/mhs294/mulhall/internals/env"
	"github.com/mhs294/mulhall/internals/types"
	"go.mongodb.org/mongo-driver/bson"
)

const teamsColl = "teams"

// TeamRepository is a mechanism for loading Team information.
// Since NFL teams rarely change, the Teams are loaded on initialization
// and then cached for the lifetime of the application.
type TeamRepository struct {
	teams map[types.TeamID]types.Team
}

// NewTeamRepository creates a new TeamRepository instance and returns a pointer to it.
// connStr is the connection string for the MongoDB instance (e.g. - "mongodb+srv://{user}:{pass}@myinstance.mongodb.net/").
func NewTeamRepository(connStr string) (*TeamRepository, error) {
	// Load Teams from database
	teams, err := loadTeams(connStr)
	if err != nil {
		return nil, err
	}

	// Initialize TeamRepository internal map
	tr := TeamRepository{teams: make(map[types.TeamID]types.Team, len(teams))}
	for _, t := range teams {
		tr.teams[t.ID] = t
	}

	return &tr, nil
}

// loadTeams loads all Teams from the database and returns them in a slice.
// If an error occurs during load, that error will be returned alongside a nil slice.
func loadTeams(connStr string) ([]types.Team, error) {
	// Create a new client and connect to the server
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, env.Timeout)
	defer cancel()
	client, err := CreateMongoDBClient(connStr, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	// Setup deferred connection closure for when function completes
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Printf("%v", err)
		}
	}()

	// Load all of the documents from the teams collection
	coll := client.Database(mulhallDB).Collection(teamsColl)
	cursor, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to query %s.%s: %v", mulhallDB, teamsColl, err)
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
func (tr *TeamRepository) GetTeam(id types.TeamID) types.Team {
	return tr.teams[id]
}
