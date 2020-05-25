package job

// A Poller polls its source(s) for the available job.
type Poller interface {
	Poll() *Request
}
