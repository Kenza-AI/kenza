package job

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"

	"github.com/kenza-ai/kenza/api/httputil"
	"github.com/kenza-ai/kenza/api/project"
	"github.com/kenza-ai/kenza/event"
	"github.com/kenza-ai/kenza/pubsub"
)

// WebhookValidator validates an incoming webhook request.
// See https://developer.github.com/webhooks/securing/ for details.
type WebhookValidator func(r *http.Request, secretToken []byte) (payload []byte, err error)

// SubmitJob submits a job on the 'arrivals' queue for the scheduler service to pick up.
func SubmitJob(jobs Store, projects project.Store, notifier pubsub.Publisher, webhookValidator WebhookValidator, webhookSecret string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close()
		payload, err := submitJobPayload(r, webhookValidator, webhookSecret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			e(err.Error())
			return
		}

		req := GitHub{}
		if err := json.Unmarshal(payload, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			e(err.Error())
			return
		}

		// If we know exactly which account and project the request is for (on-demand "run job" from the web UI)
		// we can skip traversing all projects to find all projects matching the branch/tag requested.
		projectID, accountID, ok := isRequestForSpecificProject(r)
		if ok {
			evt, err := publishJobArrivalEvent(req, projectID, accountID, notifier)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				e("failed to enqueue job '%s'", err)
				return
			}
			i("Enqueued job request %+v", evt)
			w.WriteHeader(http.StatusAccepted)
			return
		}

		// Parse delivery id
		req.DeliveryID = r.Header.Get("X-GitHub-Delivery")

		// Grab all projects watching the requested repository
		projectsWatchingRepo, err := projectsWatchingRepo(projects, req.Repo.CloneURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if len(projectsWatchingRepo) == 0 {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		i("Found projects matching repository %s: %v", req.Repo.CloneURL, projectsWatchingRepo)

		// Create a job per matching project if the ref is one we should build i.e. it matches the project's ref regex e.g. refs/heads/master
		for _, project := range projectsWatchingRepo {
			matches, err := regexp.MatchString("^"+project.Branch+"$", req.Branch)
			if !matches || err != nil {
				continue
			}

			evt, err := publishJobArrivalEvent(req, project.ID, project.AccountID, notifier)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				e("failed to enqueue job '%s'", err)
				return // TODO(ilazakis): attempt to save job, retries / DLX are two valid options
			}
			i("Enqueued job request %+v", evt)
		}

		// Respond
		w.WriteHeader(http.StatusAccepted)
	})
}

func isRequestForSpecificProject(r *http.Request) (projectID, accountID int64, ok bool) {
	accountID, err := strconv.ParseInt(httputil.Param(r, "accountID"), 10, 64)
	if err != nil {
		return -1, -1, false
	}

	projectID, err = strconv.ParseInt(httputil.Param(r, "projectID"), 10, 64)
	if err != nil {
		return -1, -1, false
	}

	return projectID, accountID, true
}

func publishJobArrivalEvent(webhook GitHub, projectID, accountID int64, notifier pubsub.Publisher) (event.JobArrived, error) {
	// Publish "job requested" event
	evt := event.JobArrived{
		AccountID:  accountID,
		CloneURL:   webhook.Repo.CloneURL,
		Ref:        webhook.Branch,
		CommitID:   webhook.Commit.ID,
		ProjectID:  projectID,
		DeliveryID: webhook.DeliveryID,
		Submitter:  webhook.Sender.Username,
	}

	return evt, notifier.Publish(evt, event.JobArrivedRoutingKey)
}

func submitJobPayload(r *http.Request, webhookValidator WebhookValidator, webhookSecret string) ([]byte, error) {
	// GitHub Webhook
	if signature := r.Header.Get("X-Hub-Signature"); signature != "" {
		return webhookValidator(r, []byte(webhookSecret))
	}

	// On-demand, from authenticated client
	return ioutil.ReadAll(r.Body)
}

func projectsWatchingRepo(projects project.Store, repo string) ([]project.Project, error) {
	return projects.GetAllWatchingRepo(repo)
}
