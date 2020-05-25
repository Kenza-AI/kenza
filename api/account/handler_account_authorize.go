package account

import (
	"net/http"

	"github.com/kenza-ai/kenza/api/httputil"
)

// Authorize validates a JWT, or API key if the caller is a Kenza service, and populates the request Context with the JWT claims.
// MUST be called before all handlers expecting authenticated requests.
func Authorize(next http.Handler, jwtSigningKey, apiKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httputil.Authorize(w, r, jwtSigningKey, apiKey, next)
	})
}
