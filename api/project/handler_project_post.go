package project

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"

	_httputil "github.com/kenza-ai/kenza/api/httputil"
)

// Create creates a project (and accompanying webhook if one is not already registered for the repo).
func Create(projects Store, client *http.Client, apiHost, webhookSecret string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accountID, err := strconv.ParseInt(_httputil.Param(r, "accountID"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			e(err.Error())
			return
		}

		type createProjectPayload struct {
			Title       string
			Description string
			Repo        string
			Branch      string
			AccessToken string
		}

		payload := createProjectPayload{}
		if err := json.Unmarshal(body, &payload); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			e(err.Error())
			return
		}

		// Create webhook
		i("Creating webhook for project at %s", payload.Repo)
		hookURL := "http://" + apiHost + "/v1/jobs/submissions" // TODO(ilazakis): SSL support
		if err := createWebHook(payload.Repo, payload.AccessToken, hookURL, webhookSecret, *client); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			e(err.Error())
			return
		}

		// Persist project
		projectID, err := projects.Create(
			accountID,
			_httputil.UserID(r),
			payload.Title,
			payload.Description,
			payload.Repo,
			payload.Branch,
			payload.AccessToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(projectID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(b)
	})
}

func createWebHook(repo, accessToken, webhookURL, secret string, client http.Client) error {

	request, err := createWebHookRequest(repo, accessToken, webhookURL, secret)
	if err != nil {
		e("error building POST webhook request %s", err)
		return err
	}

	response, err := client.Do(request)
	if err != nil {
		e("error creating webhook %s", err)
		return err
	}
	defer response.Body.Close()

	dump, err := httputil.DumpResponse(response, true)
	if err != nil {
		e(err.Error())
	}
	i("%q", dump)

	return nil
}

// https://developer.github.com/v3/repos/hooks/#create-a-hook
func createWebHookRequest(repo, accessToken, webhookURL, secret string) (*http.Request, error) {
	url := "https://api.github.com/repos/" + sanitisedRepo(repo) + "/hooks"

	type requestConfig struct {
		URL         string `json:"url"`
		ContentType string `json:"content_type"`
		Secret      string `json:"secret"`
	}

	type createRequestBody struct {
		Name   string        `json:"name"`
		Config requestConfig `json:"config"`
	}

	b := createRequestBody{
		Name: "web",
		Config: requestConfig{
			URL:         webhookURL,
			Secret:      secret,
			ContentType: "json",
		},
	}
	buffer := new(bytes.Buffer)
	if err := json.NewEncoder(buffer).Encode(b); err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", url, buffer)
	if err != nil {
		return nil, err
	}

	q := request.URL.Query()
	q.Set("access_token", accessToken)
	request.URL.RawQuery = q.Encode()

	dump, err := httputil.DumpRequest(request, true)
	if err != nil {
		e(err.Error())
	}
	i("%q", dump)

	return request, nil
}

func sanitisedRepo(repo string) string {
	repo = strings.TrimSuffix(repo, ".git")
	repo = strings.TrimPrefix(repo, "https://")
	return strings.SplitN(repo, "/", 2)[1] // drops host name, keeps GitHub account and repo e.g. kenza-ai/kenza
}
