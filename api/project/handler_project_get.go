package project

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/kenza-ai/kenza/api/httputil"
)

// GetAll returns all projects for an account.
func GetAll(projects Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accountID, err := strconv.ParseInt(httputil.Param(r, "accountID"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		projects, err := projects.GetAll(accountID)
		if err != nil {
			e(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body, err := json.Marshal(projects)
		if err != nil {
			e(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(body)
	})
}
