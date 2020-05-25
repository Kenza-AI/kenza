package api

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/kenza-ai/kenza/api/job"
	"github.com/kenza-ai/kenza/event"
)

// HTTP - Kenza HTTP `Client` implementation used by other Kenza services and the Kenza cli.
// It abstracts the underlying storage mechanism(s) and provides a universal method
// of accessing / creating / updating resources like projects, jobs, schedules etc.
type HTTP struct {
	// BaseURL of the API e.g. https://localhost:8080.
	baseURL *url.URL

	// UserAgent of caller, used to identify the caller and their version.
	userAgent string

	// The API version the client will point at.
	version string

	// HTTP client used for transport
	httpClient *http.Client

	// API key, used to authenticate with the API service
	apiKey string
}

// NewHTTPClient initializes a HTTP API `Client` for the given base URL and API version e.g. https://localhost:8080/v1
// Requires a http client for transpor, the user agent calls to the API will use and an API key.
func NewHTTPClient(baseURL, version, userAgent, apiKey string, client *http.Client) (*HTTP, error) {
	url, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &HTTP{
		apiKey:     apiKey,
		version:    version,
		baseURL:    url,
		userAgent:  userAgent,
		httpClient: client,
	}, nil
}

// AccessToken returns the vcs (e.g. GitHub) access token for the given account and project.
func (c *HTTP) AccessToken(accountID, projectID int64) (string, error) {
	path := c.version + fmt.Sprintf("/accounts/%d/projects/%d/secrets", accountID, projectID)

	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return "", err
	}

	var responseBody struct {
		AccessToken string `json:"accessToken"`
	}

	response, err := c.do(req, &responseBody)
	if err != nil {
		i("reponse: %+v", response)
		return "", err
	}
	return responseBody.AccessToken, err
}

// UpdateJob persists changes to the passed job
func (c *HTTP) UpdateJob(job event.JobUpdated) error {
	path := c.version + fmt.Sprintf("/accounts/%d/projects/%d/jobs/%d", job.AccountID, job.ProjectID, job.JobID)

	i("requesting job update: %+v", job)
	req, err := c.newRequest("PUT", path, job)
	if err != nil {
		return err
	}

	response, err := c.do(req, nil)
	if err != nil {
		i("update job response: %+v", response)
		return err
	}
	i("update job response: %+v", response)

	if response.StatusCode >= 400 {
		err = fmt.Errorf("update job error: received status code %d", response.StatusCode)
	}

	return err
}

// JobCreate creates a job
func (c *HTTP) JobCreate(accountID, projectID int64, submitter, deliveryID, revisionID string) (jobID int64, err error) {
	path := c.version + fmt.Sprintf("/accounts/%d/projects/%d/jobs", accountID, projectID)

	i("requesting job creation for account '%d' project '%d'", accountID, projectID)
	body := job.GitHub{
		DeliveryID: deliveryID,
		Commit:     job.Commit{ID: revisionID},
		Sender:     job.Sender{Username: submitter},
	}
	req, err := c.newRequest("POST", path, body)
	if err != nil {
		return -1, err
	}

	var createJobResponseBody struct {
		ID int64
	}
	response, err := c.do(req, &createJobResponseBody)
	if err != nil {
		i("create job response: %+v", response)
		return -1, err
	}

	if response.StatusCode >= 400 {
		err = fmt.Errorf("create job error: received status code %d", response.StatusCode)
	}

	i("create job response: %+v", response)
	return createJobResponseBody.ID, err
}
