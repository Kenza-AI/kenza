package account

import (
	"encoding/json"
	"net/http"

	"github.com/lib/pq"
)

type signUpRequestPayload struct {
	Email    string
	Password string
}

type signUpResponse struct {
	AccountID int64
	UserID    int64
}

// SignUp creates a new account and user
func SignUp(accounts Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload signUpRequestPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		accountID, userID, err := accounts.SignUp(payload.Email, payload.Password)
		if err != nil {
			e(err.Error())
			http.Error(w, signUpError(err).Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		response := signUpResponse{AccountID: accountID, UserID: userID}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

// https://www.postgresql.org/docs/current/static/errcodes-appendix.html
func signUpError(storeError error) error {
	if err, ok := storeError.(*pq.Error); ok {
		switch err.Code.Name() {
		case "unique_violation":
			return errAccountEmailAlreadyExists
		}
	}
	return errInternalStoreError
}
