package event

import (
	"time"
)

// JobArrived event.
//
// Publishers:
//   - api: when a job request arrives to create a job from a (web)hook, a scheduled job or the UI.
//
// Consumers:
//   - scheduler: decides if/when the job will be submitted to the worker queue for processing.
type JobArrived struct {
	ProjectID  int64
	AccountID  int64
	Ref        string
	CloneURL   string
	CommitID   string
	Submitter  string
	DeliveryID string
}

// JobInfo - JobArrived alias
type JobInfo = JobArrived

// JobQueued event.
//
// Publishers:
//   - scheduler: when a `Job arrival` is accepted for processing.
//		  Reasons for rejection: failure to verify origin, internal errors.
//
// Consumers:
//   - worker: picks up the job for processing.
type JobQueued struct {
	JobID int64
	JobInfo
}

// JobUpdated event.
//
// Publishers:
//   - worker: published multiple times throughout a job's lifecycle.
//
// Consumers:
//   - progress: updates job store accordingly.
type JobUpdated struct {
	JobID          int64
	CommitID       string
	ProjectID      int64
	AccountID      int64
	Service        string // e.g. "sagify"
	WorkDir        string
	Region         string // e.g. "us-east-1"
	Endpoint       string // endpoint name if exists
	Type           string // e.g. "training", "tuning"
	Status         string
	SageMakerJobID string
	StartedAt      time.Time
}

// Schedules holds "schedule name to schedule details" info
type Schedules map[string]*ScheduleEntry

// ScheduleEntry holds a scheduled job's details
type ScheduleEntry struct {
	When        string
	Tag         string
	Branch      string
	Description string
}

// SchedulesReceived event.
//
// Publishers:
//   - worker: published once when parsing the build file.
//
// Consumers:
//   - scheduler: updates project's schedules accordingly. Existing schedules not present
//		  in the event will be removed, existing schedules present in the event will
//                be updated. New Schedules (not already existing and added to the build file)
//                will be created.
type SchedulesReceived struct {
	AccountID int64
	ProjectID int64
	Schedules *Schedules
}

const (
	// JobArrivedRoutingKey - AMQP routing key for subscribers interested
	// in jobs arriving in the system (e.g. via a webhook).
	JobArrivedRoutingKey = "job.event.arrived"

	// JobQueuedRoutingKey - AMQP routing key for subscribers interested
	// in jobs getting queued for processing.
	JobQueuedRoutingKey = "job.event.queued"

	// JobUpdatedRoutingKey - AMQP routing key for subscribers interested
	// in a job's progress after it's been picked up for processing.
	JobUpdatedRoutingKey = "job.event.updated"

	// SchedulesReceivedRoutingKey - AMQP routing key for subscribers interested
	// in schedules in the build file parsed by workers.
	SchedulesReceivedRoutingKey = "schedules.event.received"
)
