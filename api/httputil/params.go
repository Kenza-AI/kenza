package httputil

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Param returns the string value (empty if no parameter with that name exists) of the requested path parameter.
func Param(r *http.Request, name string) string {
	return httprouter.ParamsFromContext(r.Context()).ByName(name)
}
