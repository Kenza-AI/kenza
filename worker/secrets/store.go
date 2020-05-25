package secrets

// Store provides access to secrets e.g. access tokens
type Store interface {
	// AccessToken returns the vcs (e.g. GitHub) access token for the given account and project.
	AccessToken(accountID, projectID int64) (string, error)
}
