package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mhs294/mulhall/internals/repos"
	"github.com/mhs294/mulhall/internals/types"
)

// ScheduleService represents a service for managing weekly Schedules of Matchups available for picks.
type ScheduleService struct {
	repo *repos.ScheduleRepository
}

// NewScheduleService creates a new instance of a ScheduleService and returns a pointer to it.
//
// r is the ScheduleRepository used to manage Schedule records in the database.
func NewScheduleService(r *repos.ScheduleRepository) *ScheduleService {
	return &ScheduleService{repo: r}
}

// CreateSchedule creates a new Schedule from the provided information.
// Returns an updated version of the Schedule model containing its ID after creation.
// Returns ScheduleConflictException if the request contains information that matches with another
// Schedule that already exists.
// Returns ScheduleInvalidClosesError if the request's Closes date/time falls outside the Schedule's
// calendar week.
//
// req is the CreateScheduleRequest containing the information required to create the Schedule.
//
// The Date of the CreateScheduleRequest references a week within a calendar year spanning from 03:00:00 EST/EDT
// on Tuesday to 02:59:59.999 EST/EDT the following Tuesday (this most routinely encapsulates a given week within
// an NFL season's schedule). If a Schedule already exists for the week in which the Date lies, a
// ScheduleConflictError will be returned containing the existing Schedule's unique identifier.
//
// A Date that falls at any point on a Tuesday in the Eastern US Time zone will be treated as part of the
// weekly period beginning on that Tuesday.
//
// The Year and Week of the request are purely for labelling/organizational purposes and are not used in any
// way to any specific date/time information for the Schedule being created.
func (s *ScheduleService) CreateSchedule(req *types.CreateScheduleRequest) (*types.Schedule, error) {
	// Parse the date string into a Time object.
	tz, err := time.LoadLocation("America/New_York")
	if err != nil {
		return nil, fmt.Errorf("failed to create schedule: %v", err)
	}
	date, err := time.ParseInLocation(time.DateOnly, req.Date, tz)
	if err != nil {
		return nil, fmt.Errorf("failed to create schedule: %v", err)
	}

	// If required, normalize the Schedule date
	// Treat any date that falls on a Tuesday as the starting Tuesday of a new week
	if date.In(tz).Weekday() == time.Tuesday {
		date = time.Date(date.Year(), date.Month(), date.Day(), 8, 0, 0, 0, time.UTC)
	}

	// Determine if a Schedule already exists for the provided Date
	sch, err := s.repo.GetByDateTime(date)
	if err != nil {
		return nil, err
	} else if sch != nil {
		return nil, &types.ScheduleConflictError{Date: req.Date, ID: sch.ID}
	}

	// Determine if a Schedule already exists for the provided Year/Week
	sch, err = s.repo.GetByYearAndWeek(req.Year, req.Week)
	if err != nil {
		return nil, err
	} else if sch != nil {
		return nil, &types.ScheduleConflictError{Year: req.Year, Week: req.Week, ID: sch.ID}
	}

	// Calculate the start and end date/times for the Schedule
	daysOffset := date.Weekday() - time.Tuesday
	start := time.Date(date.Year(), date.Month(), date.Day(), 3, 0, 0, 0, tz).
		Add(time.Hour * time.Duration(-24*daysOffset))
	end := start.Add(7 * 24 * time.Hour).
		Add(-1 * time.Millisecond)

	// Verify the close date/time falls between the start and end date/times
	if req.Closes.Before(start) || req.Closes.After(end) {
		return nil, &types.ScheduleInvalidClosesError{Start: start, End: end, Request: req}
	}

	// Create the Schedule
	sch = &types.Schedule{
		ID:       types.ScheduleID(uuid.NewString()),
		Year:     req.Year,
		Week:     req.Week,
		Start:    start,
		End:      end,
		Opens:    start,
		Closes:   req.Closes,
		Matchups: make(map[types.MatchupID]types.Matchup, 0),
		Active:   true,
	}

	if err = s.repo.Insert(sch); err != nil {
		return nil, err
	}

	return sch, nil
}

// GetByID gets the Schedule for the provided ID.
//
// id is the unique identifier of the Schedule to look up.
func (s *ScheduleService) GetByID(id types.ScheduleID) (*types.Schedule, error) {
	// Load the Schedule from the database
	sch, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Verify that the Schedule exists and is active
	if sch == nil || !sch.Active {
		return nil, &types.ScheduleNotFoundError{ID: id}
	}

	return sch, nil
}

// GetByYearAndWeek gets the Schedule corresponding to the provided year and week.
// Returns ScheduleNotFoundError if no Schedule exists for the provided year and week or that Schedule has been deactivated.
//
// year is the integer value of the season's starting year (e.g. - 2024 for the 2024-2025 season).
//
// week is the number of the week within the season (e.g. - 4 is the 4th week of the season).
func (s *ScheduleService) GetByYearAndWeek(year int, week int) (*types.Schedule, error) {
	// Load the Schedule from the database
	sch, err := s.repo.GetByYearAndWeek(year, week)
	if err != nil {
		return nil, err
	}

	// Verify that the Schedule exists and is active
	if sch == nil || !sch.Active {
		return nil, &types.ScheduleNotFoundError{Year: year, Week: week}
	}

	return sch, nil
}

// GetByDateTime gets the Schedule corresponding to the provided date/time.
// Returns ScheduleNotFoundError if no Schedule exists for the provided date/time or that Schedule has been deactivated.
//
// date is the [time.Time] that falls within the calendar week of the Schedule to load.
func (s *ScheduleService) GetByDateTime(date time.Time) (*types.Schedule, error) {
	// Load the Schedule from the database
	sch, err := s.repo.GetByDateTime(date)
	if err != nil {
		return nil, err
	}

	// Verify that the Schedule exists and is active
	if sch == nil || !sch.Active {
		return nil, &types.ScheduleNotFoundError{Date: date}
	}

	return sch, nil
}

// AddMatchup creates a Matchup from the provided request and adds that Matchup to the specified Schedule.
// Returns the updated state of the Schedule containing the Matchup being added.
// Returns MatchupInvalidError if the Matchup is missing required information or would not be valid within the Schedule.
//
// req is the CreateMatchupRequest containing the details of the Matchup to add to the Schedule.
func (s *ScheduleService) AddMatchup(req *types.CreateMatchupRequest) (*types.Schedule, error) {
	// Load the Schedule from the database
	sch, err := s.GetByID(req.ScheduleID)
	if err != nil {
		return nil, err
	}

	// Validate the Matchup against the current Schedule
	matchup := req.Matchup
	if err = validateMatchup(sch, matchup, ""); err != nil {
		return nil, err
	}

	// Add the Matchup to the Schedule
	id := types.MatchupID(uuid.NewString())
	sch.Matchups[id] = *matchup
	if err = s.repo.Update(sch); err != nil {
		return nil, err
	}

	return sch, nil
}

// UpdateMatchup updates an existing Matchup in a Schedule using the information in the provided request.
// Returns the updated state of the Schedule containing the Matchup being updated.
// Returns MatchupInvalidError if the Matchup is missing required information or would not be valid within the Schedule.
//
// req is the UpdateMatchupRequest containing the details of the Matchup to Update within the Schedule.
func (s *ScheduleService) UpdateMatchup(req *types.UpdateMatchupRequest) (*types.Schedule, error) {
	// Load the Schedule from the database
	sch, err := s.GetByID(req.ScheduleID)
	if err != nil {
		return nil, err
	}

	// Verify the specified Matchup exists in the Schedule
	if _, exists := sch.Matchups[req.MatchupID]; !exists {
		return nil, &types.MatchupNotFoundError{ScheduleID: req.ScheduleID, MatchupID: req.MatchupID}
	}

	// Validate the Matchup against the current Schedule
	matchup := req.Matchup
	if err = validateMatchup(sch, matchup, req.MatchupID); err != nil {
		return nil, err
	}

	// Update the Matchup in the Schedule
	sch.Matchups[req.MatchupID] = *matchup
	if err = s.repo.Update(sch); err != nil {
		return nil, err
	}

	return sch, nil
}

// RemoveMatchup removes the Matchup with the provided ID from the specified Schedule.
//
// schID is the unique identifier of the Schedule containing the Matchup to remove.
//
// mID is the unique identifier of the Matchup to remove.
func (s *ScheduleService) RemoveMatchup(schId types.ScheduleID, mID types.MatchupID) (*types.Schedule, error) {
	// Load the Schedule from the database
	sch, err := s.GetByID(schId)
	if err != nil {
		return nil, err
	}

	// Verify the specified Matchup exists in the Schedule
	if _, exists := sch.Matchups[mID]; !exists {
		return nil, &types.MatchupNotFoundError{ScheduleID: schId, MatchupID: mID}
	}

	// Remove the Matchup from the Schedule
	delete(sch.Matchups, mID)
	if err = s.repo.Update(sch); err != nil {
		return nil, err
	}

	return sch, nil
}

// Deactivate deactivates the specified Schedule (soft-delete).
//
// id is the unique identifier of the Schedule to deactivate.
func (s *ScheduleService) Deactivate(id types.ScheduleID) error {
	return s.repo.Deactivate(id)
}

func validateMatchup(sch *types.Schedule, matchup *types.Matchup, mID types.MatchupID) error {
	// Verify that the Matchup is not null
	if matchup == nil {
		return &types.MatchupInvalidError{
			ScheduleID: sch.ID,
			Matchup:    matchup,
			Reason:     "matchup is nil",
		}
	}

	// Verify that the Matchup date/time falls within the Schedule's date/time range
	if sch.Start.After(matchup.DateTime) || sch.End.Before(matchup.DateTime) {
		return &types.MatchupInvalidError{
			ScheduleID: sch.ID,
			Matchup:    matchup,
			Reason:     "matchup would fall outside of the schedule's date/time range",
		}
	}

	// Verify that Teams in Matchup aren't already featured in another Matchup on the Schedule
	for id, m := range sch.Matchups {
		if id == mID {
			continue
		}

		if m.HomeTeam == matchup.HomeTeam ||
			m.HomeTeam == matchup.AwayTeam ||
			m.AwayTeam == matchup.HomeTeam ||
			m.AwayTeam == matchup.AwayTeam {
			return &types.MatchupInvalidError{
				ScheduleID: sch.ID,
				Matchup:    matchup,
				Reason:     "matchup would create a conflict with teams from another existing matchup",
			}
		}
	}

	return nil
}
