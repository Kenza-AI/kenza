package job

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/kenza-ai/kenza/api/httputil"
	"github.com/kenza-ai/kenza/api/project"
)

// CreateJob creates a job (not to be confused with submitting a job which happens earlier e.g. from GitHub).
// For CreateJob to be called, a job submission has been authorized and submitted to the scheduler service which
// decides when to actually create a job entry in the database (this call).
func CreateJob(jobs Store, projects project.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		accountID, err := strconv.ParseInt(httputil.Param(r, "accountID"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			e(err.Error())
		}

		projectID, err := strconv.ParseInt(httputil.Param(r, "projectID"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			e(err.Error())
		}

		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			e(err.Error())
		}

		req := GitHub{}
		if err := json.Unmarshal(payload, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			e(err.Error())
			return
		}

		// Parse delivery id
		req.DeliveryID = r.Header.Get("X-GitHub-Delivery")

		// Create job
		jobID, err := jobs.CreateJob(accountID, projectID, req.Sender.Username, req.DeliveryID, req.Commit.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			e("error creating job: %s", err.Error())
			return
		}
		i(`Created job "%d" in account "%d" in project "%d"`, jobID, accountID, projectID)

		var createJobResponseBody = struct {
			ID int64
		}{
			ID: jobID,
		}
		body, err := json.Marshal(createJobResponseBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(body)
	})
}
