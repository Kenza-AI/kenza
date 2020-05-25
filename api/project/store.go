package project

// The Project domain model.
type Project struct {
	ID          int64  `json:"projectID"`
	AccountID   int64  `json:"accountID"`
	CreatorID   int64  `json:"creatorID"`
	Creator     string `json:"creator"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Repo        string `json:"repo"`
	Branch      string `json:"branch"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
}

// Store abstracts communication with the service's persistance mechanism(s).
type Store interface {
	// Delete deletes a project.
	Delete(accountID, projectID int64) error

	// GetAccessToken returns the vcs (e.g. GitHub) access token for the given account and project.
	GetAccessToken(accountID, projectID int64) (string, error)

	// GetAll returns all projects under the account ID provided.
	GetAll(accountID int64) ([]Project, error)

	// GetAllWatchingRepo returns all projects watching for changes to the repo provided.
	GetAllWatchingRepo(repo string) ([]Project, error)

	// ReferencesForRepo returns the allowed references regex (branches and tags a project should build) for a repository.
	// The returnd result is a map of project ids to allowed regex.
	ReferencesForRepo(repository string) (regexPerProject map[int64]string, err error)

	// Create creates a new project.
	Create(accountID, creatorID int64, title, description, repoURL, branch, githubAccessToken string) (projectID int64, err error)
}
