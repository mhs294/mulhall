package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mhs294/mulhall/internals/repos"
	"github.com/mhs294/mulhall/internals/types"
)

// ScheduleService represents a service for managing weekly Schedules of Matchups available for Picks.
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
		Matchups: make([]types.Matchup, 0),
		Active:   true,
	}

	if err = s.repo.Insert(sch); err != nil {
		return nil, err
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

// TODO - AddMatchup

// TODO - EditMatchup

// TODO - RemoveMatchup

// TODO - Deactivate
