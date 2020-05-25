package job

// Notifier — notifies about the progress/status of a running job.
type Notifier interface {
	Notify() error
	NotifySchedules() error
}
