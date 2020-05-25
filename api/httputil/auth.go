package httputil

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type contextKey string

var userContextKey = contextKey("user")

var (
	errInvalidAuthHelperUsage = errors.New("called UserID before Authorize Handler or on non-authenticated endpoint")
)

type claims struct {
	Accounts map[int64]string `json:"accounts"`
	Username string           `json:"username"`
	jwt.StandardClaims
}

// Authorize validates the JWT, or API key if caller is a Kenza service, and makes claims available for downstream handlers.
func Authorize(w http.ResponseWriter, r *http.Request, jwtSigningKey, expectedAPIKey string, next http.Handler) {
	jwtString, apiKey := r.Header.Get("Authorization"), r.Header.Get("X-API-Key")

	// If no JWT was provided this is either:
	//
	// 	1. A Kenza service calling, in which case access is determined by the API key.
	//	   All services are provided with an API key for service-to-service comms.
	//
	// 	2. A VCS e.g. GitHub webhook, in which case access is determined by validating the 'X-Hub-Signature'.
	//	   See https://developer.github.com/webhooks/securing/ (webhook vaidation currrently in the "create job" handler).
	//
	// 	3. Unauthorized
	if jwtString == "" {
		if apiKey == expectedAPIKey {
			next.ServeHTTP(w, r)
			return // service-to-service comms, no claims to set, move to next handler.
		}
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	claims := &claims{}
	tkn, err := jwt.ParseWithClaims(jwtString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSigningKey), nil
	})

	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	if !tkn.Valid {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Add claims to the request's context for handlers down the chain to be
	// able to use helpers e.g. `UserID(r *http.Request`) to retrieve auth/user info.
	ctx := context.WithValue(r.Context(), userContextKey, claims)
	next.ServeHTTP(w, r.WithContext(ctx))
}

// JWT signs a new JWT with the parameters as claims (plus some standard claims).
func JWT(username string, userID int64, accounts map[int64]string, jwtKey string) (token string, err error) {
	claims := &claims{
		Accounts: accounts,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			Subject:   strconv.FormatInt(userID, 10),
			ExpiresAt: time.Now().Add(60 * time.Minute).Unix(),
		},
	}

	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tkn.SignedString([]byte(jwtKey))
}

// UserID returns the authenticated user's ID for the request.
// MUST be called after the `Authorize` handler, panics otherwise.
func UserID(r *http.Request) int64 {
	claims, ok := r.Context().Value(userContextKey).(*claims)
	if !ok {
		panic(errInvalidAuthHelperUsage)
	}

	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		panic(errInvalidAuthHelperUsage)
	}
	return userID
}

func (c contextKey) String(s string) string {
	return "ai.kenza.context.key." + string(c)
}
