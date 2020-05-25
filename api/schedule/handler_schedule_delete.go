package schedule

import (
	"net/http"
	"strconv"

	"github.com/kenza-ai/kenza/api/httputil"
)

// Delete a schedule.
func Delete(schedules Store) http.Handler {
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

		scheduleID, err := strconv.ParseInt(httputil.Param(r, "scheduleID"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			e(err.Error())
		}

		// Delete schedule
		if err = schedules.Delete(accountID, projectID, scheduleID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			e("error deleting schedule: %s", err.Error())
			return
		}

		i(`Deleted schedule "%d" for project "%d" in account "%d"`, scheduleID, projectID, accountID)
		w.WriteHeader(http.StatusOK)
	})
}
