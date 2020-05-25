package schedule

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/kenza-ai/kenza/api/httputil"
)

// Update a schedule.
func Update(schedules Store) http.Handler {
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

		// Update schedule
		scheduleID, err := schedules.Update(schedule)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			e("error creating schedule: %s", err.Error())
			return
		}
		i(`Updated schedule "%d" for project "%d" in account "%d"`, scheduleID, projectID, accountID)

		var updateScheduleResponseBody = struct {
			ID int64
		}{
			ID: schedule.ID,
		}
		body, err := json.Marshal(updateScheduleResponseBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(body)
	})
}
