package handler

import (
	"log"
	"path/filepath"
	"strconv"
	"time"

	"github.com/kenza-ai/kenza/worker/job"
)

// The Arrival Handler is the first handler called after a job is received for processing.
type Arrival struct {
	workDir string
	next    job.Handler
}

// NewArrival — Arrival Handler initialiser
func NewArrival(workDir string) *Arrival {
	return &Arrival{workDir: workDir}
}

// Handle implementation for the Arrival Handler.
//
// 1. Sets job status to running / in progress
// 2. Sets job work / cloning directory
// 3. Notifies about the new status and clone directory
func (h *Arrival) Handle(r *job.Request) {
	r.Status = "running"
	r.StartedAt = time.Now().UTC()
	r.WorkDir = filepath.Join(h.workDir, strconv.FormatInt(r.JobQueued.ProjectID, 10), strconv.FormatInt(r.JobQueued.JobID, 10))

	if err := r.Notify(); err != nil {
		log.Print(err)
	}

	if h.next != nil {
		h.next.Handle(r)
	}
}

// SetNext sets the next Handler
func (h *Arrival) SetNext(next job.Handler) {
	h.next = next
}
