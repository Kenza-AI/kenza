package project

import (
	"database/sql"
)

// Postgres is a postgres projects `Store` implementation.
type Postgres struct {
	DB *sql.DB
}

// GetAll returns all projects under the account ID provided.
func (store *Postgres) GetAll(accountID int64) ([]Project, error) {
	rows, err := store.DB.Query(getProjectsStatement, accountID)
	if err != nil {
		return []Project{}, err
	}

	projects := []Project{}
	for rows.Next() {
		var project Project
		err = rows.Scan(
			&project.ID,
			&project.AccountID,
			&project.CreatorID,
			&project.Title,
			&project.Description,
			&project.Repo,
			&project.Branch,
			&project.Created,
			&project.Updated,
			&project.Creator)

		if err != nil {
			return []Project{}, err
		}
		projects = append(projects, project)
	}

	return projects, err
}

// GetAllWatchingRepo returns all projects watching for chanegs to the repo provided.
func (store *Postgres) GetAllWatchingRepo(repo string) ([]Project, error) {
	rows, err := store.DB.Query(getProjectsWatchingRepoStatement, repo)
	if err != nil {
		return []Project{}, err
	}

	projects := []Project{}
	for rows.Next() {
		var project Project
		err = rows.Scan(
			&project.ID,
			&project.AccountID,
			&project.CreatorID,
			&project.Title,
			&project.Description,
			&project.Repo,
			&project.Branch,
			&project.Created,
			&project.Updated,
			&project.Creator)

		if err != nil {
			return []Project{}, err
		}
		projects = append(projects, project)
	}

	return projects, err
}

// ReferencesForRepo returns the allowed references regex (branches and tags a project should build) for a repository.
// The returnd result is a map of project ids to allowed regex.
func (store *Postgres) ReferencesForRepo(repository string) (regexPerProject map[int64]string, err error) {
	regexPerProject = map[int64]string{}
	rows, err := store.DB.Query(getRefsForRepoStatement, repository)
	if err != nil {
		return regexPerProject, err
	}

	for rows.Next() {
		var projectID int64 = -1
		refsRegex := ""
		err = rows.Scan(&projectID, &refsRegex)

		if err != nil {
			return regexPerProject, err
		}
		regexPerProject[projectID] = refsRegex
	}

	return regexPerProject, err
}

// Create creates a new project.
func (store *Postgres) Create(accountID, creatorID int64, title, description, repo, refs, vcsAccessToken string) (projectID int64, err error) {
	err = store.DB.QueryRow(createProjectStatement, accountID, creatorID, title, description, repo, refs, vcsAccessToken).Scan(&projectID)
	return projectID, err
}

// GetAccessToken returns the vcs (e.g. GitHub) access token for the given account and project.
func (store *Postgres) GetAccessToken(accountID, projectID int64) (accessToken string, err error) {
	err = store.DB.QueryRow(getAccessTokenStatement, accountID, projectID).Scan(&accessToken)
	return accessToken, err
}

// Delete a project.
func (store *Postgres) Delete(accountID, projectID int64) error {
	const deleteProjectStatement = `DELETE FROM projects WHERE id = $1 AND account_id = $2 RETURNING id`

	var id int64
	return store.DB.QueryRow(deleteProjectStatement, projectID, accountID).Scan(&id)
}

const getAccessTokenStatement = `
SELECT vcs_access_token
FROM kenza.projects
WHERE account_id = $1 AND id = $2
`

const getProjectsStatement = `
SELECT projects.id, account_id, creator_id, title, description, repository, refs, projects.created, projects.updated, users.username
FROM kenza.projects
INNER JOIN kenza.users ON users.id = projects.creator_id
WHERE account_id = $1
`

const getProjectsWatchingRepoStatement = `
SELECT projects.id, account_id, creator_id, title, description, repository, refs, projects.created, projects.updated, users.username
FROM kenza.projects
INNER JOIN kenza.users ON users.id = projects.creator_id
WHERE LOWER(repository) = LOWER($1)
`

const getRefsForRepoStatement = `
SELECT id, refs
FROM kenza.projects
WHERE LOWER(repository) = LOWER($1)
`

const createProjectStatement = `
INSERT INTO kenza.projects (account_id, creator_id, title, description, repository, refs, vcs_access_token)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id
`
