package job

import (
	"net/http"
	"strconv"

	"github.com/kenza-ai/kenza/api/httputil"
)

// Cancel jobs
func Cancel(jobs Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accountID, err := strconv.ParseInt(httputil.Param(r, "accountID"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		projectID, err := strconv.ParseInt(httputil.Param(r, "projectID"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		jobsToCancelStr, _ := r.URL.Query()["id"]
		jobIDsToCancel := []int64{}
		for _, job := range jobsToCancelStr {
			jobID, err := strconv.ParseInt(job, 10, 64)
			if err != nil {
				e("error parsing job id query param into integer: %s", err.Error())
				continue
			}
			jobIDsToCancel = append(jobIDsToCancel, jobID)
		}

		err = jobs.CancelJobs(accountID, projectID, jobIDsToCancel)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
