package schedule

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/kenza-ai/kenza/api/httputil"
)

// Create a schedule that will run a job based on a cron expression.
func Create(schedules Store) http.Handler {
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

		schedule := Schedule{}
		if err := json.Unmarshal(payload, &schedule); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			e(err.Error())
			return
		}

		// Create schedule
		scheduleID, err := schedules.Create(schedule)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			e("error creating schedule: %s", err.Error())
			return
		}
		i(`Created schedule "%d" for project "%d" in account "%d"`, scheduleID, projectID, accountID)

		var createScheduleResponseBody = struct {
			ID int64
		}{
			ID: scheduleID,
		}
		body, err := json.Marshal(createScheduleResponseBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(body)
	})
}
