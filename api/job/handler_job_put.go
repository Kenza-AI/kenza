package job

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/kenza-ai/kenza/event"
)

// Put updates the job's details based on `event.JobUpdated`
func Put(jobs Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			e(err.Error())
			return
		}

		var update event.JobUpdated
		if err := json.Unmarshal(body, &update); err != nil {
			e(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := jobs.UpdateJob(update); err != nil {
			e(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
