package api

import (
	"github.com/kenza-ai/kenza/api/schedule"
	"github.com/kenza-ai/kenza/event"
)

// Client is the Kenza API client used by other Kenza services and the Kenza cli.
// It abstracts the underlying storage mechanism(s) and provides a universal method
// of accessing / creating / updating resources like projects, jobs, schedules etc.
type Client interface {
	Jobs
	Projects
	Schedules
}

// Jobs client
type Jobs interface {
	// UpdateJob persists changes to the passed job.
	UpdateJob(job event.JobUpdated) error

	// JobCreate - creates a job.
	JobCreate(accountID, projectID int64, submitter, deliveryID, revisionID string) (jobID int64, err error)
}

// Projects client
type Projects interface {
	// AccessToken returns the vcs (e.g. GitHub) access token for the given account and project.
	AccessToken(accountID, projectID int64) (string, error)
}

// Schedules client
type Schedules interface {
	// Schedules returns all schedules, excluding any ids in the excluded list.
	Schedules(excludedIDs []int64) ([]schedule.Schedule, error)

	// SchedulesForProject returns a project's schedules.
	SchedulesForProject(accountID, projectID int64) ([]schedule.Schedule, error)

	// Create a schedule.
	ScheduleCreate(schedule schedule.Schedule) (scheduleID int64, err error)

	// Update a schedule.
	ScheduleUpdate(schedule schedule.Schedule) error

	// Delete a schedule.
	ScheduleDelete(accountID, projectID, scheduleID int64) error
}
