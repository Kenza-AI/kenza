package vcs

import (
	"errors"
	"log"
	"os"

	gitvcs "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

type repository struct {
	URL    string `json:"url"`
	Branch string `json:"branch"`
}

// Git provides a VCS interface implementation for git repositories.
type Git struct{}

// Scheduled jobs have no commit ID attached to them
// since there's no way to predict a future commit ID
const scheduledJobCommitID = "scheduledjob"

var (
	errCommitDoesNotMatchHeadAfterCheckout = errors.New("HEAD not matching requested commit after git checkout. This should never happen")
)

// Checkout clones the repo at "repoURL" and checks out "commit" — or HEAD if no commit is
// provided — of "branch" to the specified path. The "path" directory (and subdirectories)
// will be created if needed. A GitHub access token is required for private repos.
func (git *Git) Checkout(repoURL, branch, commit, path, gitHubAccessToken string) (commitID string, err error) {

	if _, err := os.Stat(path); err == nil {
		log.Printf("Git repository at %q already exists, removing all content to proceed", path)
		if err := os.RemoveAll(path); err != nil {
			return "", err
		}
	}

	gitCloneOptions := &gitvcs.CloneOptions{
		URL: repoURL,
		// If the requested commit is not in the shallow repo after cloning we return an error.
		// We can clone all history to avoid such cases (remove depth) or clone in a loop
		// (with a maximum number of attempts to avoid potential "infinite" loops) until we find the commit.
		Depth:         15,
		ReferenceName: plumbing.ReferenceName(branch),
		SingleBranch:  true,
		Progress:      os.Stdout,
		Tags:          gitvcs.NoTags,
	}

	if gitHubAccessToken != "" {
		auth := &http.BasicAuth{Username: "any username will do when using tokens", Password: gitHubAccessToken}
		gitCloneOptions.Auth = auth
	}

	repo, err := gitvcs.PlainClone(path, false, gitCloneOptions)
	if err != nil {
		return "", err
	}

	// When no specific commit is requested as is the case
	// with 'scheduled for later' jobs we simply checkout
	// whichever commit is HEAD at this point in time.
	if len(commit) < 1 || commit == scheduledJobCommitID {
		return checkoutHead(repo)
	}

	return checkoutCommit(commit, repo)
}

func checkoutHead(repo *gitvcs.Repository) (commitID string, err error) {
	head, err := repo.Head() // HEAD is implicitly checked out post-cloning in git-go
	if err != nil {
		return "", err
	}
	return head.Hash().String(), nil
}

func checkoutCommit(commit string, repo *gitvcs.Repository) (commitID string, err error) {
	// We have been asked to check a specific commit out (commit is not empty)
	worktree, err := repo.Worktree()
	if err != nil {
		return "", err
	}

	if err = worktree.Checkout(&gitvcs.CheckoutOptions{
		Hash: plumbing.NewHash(commit),
	}); err != nil {
		return "", err
	}

	// Sanity check that HEAD now points to 'commit'
	head, err := repo.Head()
	if err != nil {
		return "", err
	}

	if head.Hash().String() != commit {
		return "", errCommitDoesNotMatchHeadAfterCheckout
	}

	return commit, nil
}
