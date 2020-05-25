package schedule

// Schedule - defines how often a job will run against the specified repository and VCS ref e.g. a specific branch
type Schedule struct {
	ID          int64
	ProjectID   int64
	AccountID   int64
	Title       string
	Cron        string
	Repository  string
	Description string
	Ref         string // e.g. refs/heads/master or refs/tags/v1.7.3
}

// Store - schedule store abstraction
type Store interface {
	// GetAll all schedules, excluding any ids in the excluded list.
	GetAll(excludedIDs []int64) ([]Schedule, error)

	// GetSchedulesForProject returns a project's schedules.
	GetSchedulesForProject(accountID, projectID int64) ([]Schedule, error)

	// Create a schedule.
	Create(schedule Schedule) (scheduleID int64, err error)

	// Update a schedule.
	Update(schedule Schedule) (scheduleID int64, err error)

	// Delete a schedule.
	Delete(accountID, projectID, scheduleID int64) error
}
