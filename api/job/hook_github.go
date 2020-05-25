package job

// GitHub – incoming GitHub webhook payload. Only exporting fields Kenza needs.
// https://developer.github.com/v3/activity/events/types/#pushevent
type GitHub struct {
	DeliveryID string
	Branch     string     `json:"ref"`
	Sender     Sender     `json:"sender"`
	Repo       repository `json:"repository"`
	Commit     Commit     `json:"head_commit"`
}

type repository struct {
	CloneURL string `json:"clone_url"`
}

type Sender struct {
	Username string `json:"login"`
}

type Commit struct {
	ID string `json:"id"`
}
