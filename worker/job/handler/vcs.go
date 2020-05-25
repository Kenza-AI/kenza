package handler

import (
	"github.com/kenza-ai/kenza/worker/job"
	"github.com/kenza-ai/kenza/worker/secrets"
	versioncontrol "github.com/kenza-ai/kenza/worker/vcs"
)

// The VCS Handler clones the repo and returns the checked–out commit ID.
type VCS struct {
	sec       secrets.Store
	vcsClient versioncontrol.Client
	next      job.Handler
}

// NewVCS — VCS Handler constructor
func NewVCS(versionControl versioncontrol.Client, sec secrets.Store) *VCS {
	return &VCS{vcsClient: versionControl, sec: sec}
}

// Handle implementation for the VCS Handler.
//
// 1. Clones repo
// 2. Checks out commit ID or HEAD if no commit ID is provided
// 3. Notifies about the commit ID
// 4. Calls next handler if one is provided
func (h *VCS) Handle(r *job.Request) {

	accessToken, err := h.sec.AccessToken(r.JobQueued.AccountID, r.JobQueued.ProjectID)
	if err != nil {
		r.Fail(err)
		return
	}

	checkedOutCommit, err := h.vcsClient.Checkout(r.CloneURL, r.Ref, r.JobQueued.CommitID, r.WorkDir, accessToken)
	if err != nil {
		r.Fail(err)
		return
	}
	i("cloned repo %s and commit id %s in %s", r.CloneURL, checkedOutCommit, r.WorkDir)
	r.JobUpdated.CommitID = checkedOutCommit

	if err := r.Notify(); err != nil {
		e(err.Error())
	}

	if h.next != nil {
		h.next.Handle(r)
	}
}

// SetNext sets the next Handler
func (h *VCS) SetNext(next job.Handler) {
	h.next = next
}
