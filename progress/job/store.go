package job

import "github.com/kenza-ai/kenza/event"

// Store abstracts the underlying persistence mechanism for Job updates
type Store interface {
	// UpdateJob persists changes to the passed job
	UpdateJob(job event.JobUpdated) error
}
