package account

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/kenza-ai/kenza/api/httputil"
)

// SignIn creates and returns a new access token for the user
func SignIn(accounts Store, jwtKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		type credentials struct {
			Email    string
			Password string
		}

		var creds credentials
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		userID, accounts, err := accounts.SignIn(creds.Email, creds.Password)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		tokenString, err := httputil.JWT(creds.Email, userID, accounts, jwtKey)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		type response struct {
			UserID   string           `json:"userID"`
			Accounts map[int64]string `json:"accounts"`
			Username string           `json:"username"`
			JWT      string           `json:"accessToken"`
		}

		resp := &response{JWT: tokenString, Accounts: accounts, Username: creds.Email, UserID: strconv.FormatInt(userID, 10)}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
