package secrets

import (
	"net/http"
	"time"

	"github.com/kenza-ai/kenza/api/api"
)

// NewSecretsClient wraps a `client.HTTP` API `Client` and forwards requests for VCS access
// tokens (secrets) to the Kenza API (as opposed to coupling a store / DB to the worker).
func NewSecretsClient(apiKey, apiVersion, workerVersion string) (*api.HTTP, error) {
	httpClient := &http.Client{
		Timeout: time.Second * 15,
	}
	userAgent := "Worker/" + workerVersion

	return api.NewHTTPClient("http://api:8080", apiVersion, userAgent, apiKey, httpClient)
}
