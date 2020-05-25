package project

import (
	"net/http"
	"strconv"

	"github.com/kenza-ai/kenza/api/httputil"
)

// Delete a project.
func Delete(projects Store) http.Handler {
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

		if err := projects.Delete(accountID, projectID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			e("error deleting project: %s", err.Error())
			return
		}

		i(`Deleted project "%d" in account "%d"`, projectID, accountID)
		w.WriteHeader(http.StatusOK)
	})
}
