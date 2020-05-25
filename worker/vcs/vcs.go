package vcs

// Client provides repo related functionality (e.g cloning repos)
type Client interface {
	// Checkout clones the repo at "repoURL" and checks out "commit" — or HEAD if no commit is
	// provided — of "branch" to the specified path. The "path" directory (and subdirectories)
	// will be created if needed. A GitHub access token is required for private repos.
	Checkout(repoURL, branch, commit, path, gitHubAccessToken string) (commitSHA string, err error)
}
