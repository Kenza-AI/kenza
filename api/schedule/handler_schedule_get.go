package schedule

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/kenza-ai/kenza/api/httputil"
)

// GetSchedulesForProject returns a project's schedules.
func GetSchedulesForProject(schedules Store) http.Handler {
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

		schedulesForProject, err := schedules.GetSchedulesForProject(accountID, projectID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			e("error retrieving schedules: %s", err.Error())
			return
		}
		i(`retrieved schedules "%+v" for project "%d" in account "%d"`, schedulesForProject, projectID, accountID)

		var schedulesForProjectResponseBody = struct {
			Schedules []Schedule
		}{
			Schedules: schedulesForProject,
		}
		body, err := json.Marshal(schedulesForProjectResponseBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(body)
	})
}

// GetSchedules returns all schedules outside the ones matching the exclusion list of schedule ids.
func GetSchedules(schedules Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		excludedSchedulesStr, _ := r.URL.Query()["excluded"]
		excludedSchedules := []int64{}
		for _, schedule := range excludedSchedulesStr {
			scheduleID, err := strconv.ParseInt(schedule, 10, 64)
			if err != nil {
				e("error parsing schedule id query param into integer: %s", err.Error())
				continue
			}
			excludedSchedules = append(excludedSchedules, scheduleID)
		}

		allSchedules, err := schedules.GetAll(excludedSchedules)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			e("error retrieving schedules: %s", err.Error())
			return
		}
		i(`retrieved schedules "%+v"`, allSchedules)

		var schedulesResponseBody = struct {
			Schedules []Schedule
		}{
			Schedules: allSchedules,
		}
		body, err := json.Marshal(schedulesResponseBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(body)
	})
}
