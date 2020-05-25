package job

import (
	"net/http"
	"time"

	"github.com/kenza-ai/kenza/api/api"
)

// NewJobsClient wraps a `client.HTTP` API `Client` and forwards "job updated"
// requests to the Kenza API (as opposed to directly depending on a store / DB).
func NewJobsClient(apiKey, apiVersion, progressVersion string) (*api.HTTP, error) {
	httpClient := &http.Client{
		Timeout: time.Second * 15,
	}
	userAgent := "Progress/" + progressVersion

	return api.NewHTTPClient("http://api:8080", apiVersion, userAgent, apiKey, httpClient)
}
