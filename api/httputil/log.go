package httputil

import (
	"net/http"
	"net/http/httputil"
)

// Log logs the incoming request and calls `ServeHTTP` on the next handler.
func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			e(err.Error())
		}
		i("incoming request: \n%s\n", string(dump))
		next.ServeHTTP(w, r)
	})
}
