package repos

import (
	"fmt"
	"time"

	"github.com/mhs294/mulhall/internals/db"
	"github.com/mhs294/mulhall/internals/types"
	"go.mongodb.org/mongo-driver/bson"
)

// ScheduleRepository manages Schedule records in the database.
type ScheduleRepository struct {
	mdb      *db.MongoDB
	dbName   string
	collName string
}

// NewScheduleRepository creates a new ScheduleRepository instance and returns a pointer to it.
//
// mdb is the MongoDB instance used by the ScheduleRepository.
func NewScheduleRepository(mdb *db.MongoDB) *ScheduleRepository {
	return &ScheduleRepository{mdb: mdb, dbName: "mulhall", collName: "schedules"}
}

// TestConnection tests the connection to the MongoDB instance and returns any error that occurs.
func (r *ScheduleRepository) TestConnection() error {
	return r.mdb.TestConnection(r.dbName)
}

// Insert inserts the provided Schedule into the database.
//
// s is the Schedule to insert into the database.
func (r *ScheduleRepository) Insert(s *types.Schedule) error {
	if err := r.mdb.InsertOne(r.dbName, r.collName, s); err != nil {
		return fmt.Errorf("failed to insert schedule: %v", err)
	}

	return nil
}

// GetByDateTime gets the Schdule whose start/end window contains the provided date/time.
//
// dateTime is the [time.Time] whose corresponding Schedule will be loaded.
func (r *ScheduleRepository) GetByDateTime(datetime *time.Time) (*types.Schedule, error) {
	// Define the query
	query := bson.M{
		"start": bson.M{"$lte": datetime},
		"end":   bson.M{"$gte": datetime},
	}

	// Load the Schedule from the database
	var s types.Schedule
	if err := r.mdb.GetOne(r.dbName, r.collName, query, &s); err != nil {
		return nil, fmt.Errorf("failed to look up schedule (datetime=%s): %v", datetime.Format(time.UnixDate), err)
	}

	return &s, nil
}

// GetByYearAndWeek gets the Schedule corresponding to the provided year and week of a season.
//
// year is the integer value of the season's starting year (e.g. - 2024 for the 2024-2025 season).
//
// week is the number of the week within the season (e.g. - 4 is the 4th week of the season).
func (r *ScheduleRepository) GetByYearAndWeek(year int, week int) (*types.Schedule, error) {
	// Define the query
	query := bson.M{
		"year": year,
		"week": week,
	}

	// Load the Schedule from the database
	var s types.Schedule
	if err := r.mdb.GetOne(r.dbName, r.collName, query, &s); err != nil {
		return nil, fmt.Errorf("failed to look up schedule (year=%d, week=%d): %v", year, week, err)
	}

	return &s, nil
}

// Update updates a Schedule in the database using the information in the provided model.
//
// e is the model to use to update the Schedule in the database. The models' ScheduleID is used to
// determine which document in the database should be replaced with the updated version.
func (r *ScheduleRepository) Update(s *types.Schedule) error {
	// Define the filter query
	filter := bson.M{"id": s.ID}

	// Perform the update
	if err := r.mdb.ReplaceOne(r.dbName, r.collName, filter, s); err != nil {
		return fmt.Errorf("failed to update schedule (%v): %v", s, err)
	}

	return nil
}
