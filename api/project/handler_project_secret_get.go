package project

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/kenza-ai/kenza/api/httputil"
)

// GetAccessToken returns the access token (e.g. GitHub) for the project.
func GetAccessToken(projects Store) http.Handler {
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

		accessToken, err := projects.GetAccessToken(accountID, projectID)
		if err != nil {
			e(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		responseBody := struct {
			AccessToken string `json:"accessToken"`
		}{
			AccessToken: accessToken,
		}

		body, err := json.Marshal(responseBody)
		if err != nil {
			e(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(body)
	})
}
