package job

import (
	"database/sql"
	"time"

	"github.com/kenza-ai/kenza/event"
	"github.com/lib/pq"
)

// Postgres is the postgres jobs `Store` implementation.
type Postgres struct {
	DB *sql.DB
}

// CreateJob creates a new training job.
func (store *Postgres) CreateJob(accountID, projectID int64, submitter, deliveryID, revisionID string) (jobID int64, err error) {
	err = store.DB.QueryRow(
		createJobStatement,
		projectID,
		delivery(deliveryID),
		"submitted",
		submitter,
		revisionID).Scan(&jobID)
	return jobID, err
}

// Get returns a job's info.
func (store *Postgres) Get(jobID int64) (Job, error) {
	row := store.DB.QueryRow(getJobStatement, jobID)

	var job Job
	var startedValue sql.NullTime                                   // nullable columns
	var sagemakerIDValue, endpointValue, regionValue sql.NullString // nullable columns
	err := row.Scan(&job.ID, &job.Status, &job.Submitter, &job.CommitID, &sagemakerIDValue, &job.Type, &regionValue, &endpointValue, &job.Created, &job.Updated, &startedValue)
	if err != nil {
		return Job{}, err
	}

	if startedValue.Valid {
		job.Started = startedValue.Time
	}

	if regionValue.Valid {
		job.Region = regionValue.String
	}

	if endpointValue.Valid {
		job.Endpoint = endpointValue.String
	}

	if sagemakerIDValue.Valid {
		job.SageMakerID = sagemakerIDValue.String
	}

	return job, nil
}

// GetAll returns a project's jobs.
func (store *Postgres) GetAll(accountID, projectID int64) ([]Job, error) {
	rows, err := store.DB.Query(getJobsStatement, accountID, projectID)
	if err != nil {
		return []Job{}, err
	}

	jobs := []Job{}
	for rows.Next() {
		var job Job
		var startedValue sql.NullTime                                   // nullable columns
		var sagemakerIDValue, endpointValue, regionValue sql.NullString // nullable columns
		err := rows.Scan(
			&job.ID,
			&job.Status,
			&job.Submitter,
			&job.CommitID,
			&job.Created,
			&job.Updated,
			&startedValue,
			&sagemakerIDValue,
			&job.Type,
			&regionValue,
			&endpointValue,
			&job.Project.ID,
			&job.Project.Name,
			&job.Project.Repo)
		if err != nil {
			return []Job{}, err
		}

		if startedValue.Valid {
			job.Started = startedValue.Time
		}

		if regionValue.Valid {
			job.Region = regionValue.String
		}

		if endpointValue.Valid {
			job.Endpoint = endpointValue.String
		}

		if sagemakerIDValue.Valid {
			job.SageMakerID = sagemakerIDValue.String
		}

		jobs = append(jobs, job)
	}
	return jobs, err
}

// UpdateJob persists changes to the passed job.
func (store *Postgres) UpdateJob(job event.JobUpdated) error {
	var id int64
	return store.DB.QueryRow(
		updateJobStatement,
		job.JobID,
		job.CommitID,
		sagemakerID(job),
		job.Status,
		jobType(job),
		endpoint(job),
		job.Region,
		startedAt(job.StartedAt)).Scan(&id)
}

// DeleteJobs delets a list of jobs.
func (store *Postgres) DeleteJobs(accountID, projectID int64, jobIDs []int64) error {
	_, err := store.DB.Query(deleteJobsStatement, projectID, pq.Array(jobIDs))
	return err
}

// CancelJobs delets a list of jobs.
func (store *Postgres) CancelJobs(accountID, projectID int64, jobIDs []int64) error {
	_, err := store.DB.Query(cancelJobsStatement, projectID, pq.Array(jobIDs))
	return err
}

func jobType(job event.JobUpdated) string {
	if job.Type == "" {
		return "unknown"
	}
	return job.Type
}

func sagemakerID(job event.JobUpdated) sql.NullString {
	return sql.NullString{String: job.SageMakerJobID, Valid: job.Service == "sagify" && job.SageMakerJobID != ""}
}

func endpoint(job event.JobUpdated) sql.NullString {
	return sql.NullString{String: job.Endpoint, Valid: job.Endpoint != ""}
}

func delivery(deliveryID string) sql.NullString {
	return sql.NullString{String: deliveryID, Valid: deliveryID != ""}
}

func startedAt(startedAt time.Time) sql.NullTime {
	return sql.NullTime{Time: startedAt, Valid: !startedAt.IsZero()}
}

const createJobStatement = `
INSERT INTO jobs (project_id, delivery_id, status, submitter, commit_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id
`

const getJobStatement = `
SELECT id, status, submitter, commit_id, sagemaker_id, type, region, endpoint, created, updated, started
FROM kenza.jobs
WHERE id = $1
`

const getJobsStatement = `
SELECT jobs.id, status, submitter, commit_id, jobs.created, jobs.updated, jobs.started, sagemaker_id, type, region, endpoint, project_id, title, repository
FROM kenza.jobs
INNER JOIN kenza.projects ON jobs.project_id = projects.id
WHERE projects.account_id = $1 AND project_id = $2
`

const updateJobStatement = `
UPDATE kenza.jobs
SET commit_id = $2, sagemaker_id = $3, type = $5, endpoint = $6, region = $7, started = $8,
status = (
	CASE
	WHEN status = 'cancelled' 
	THEN status
	ELSE $4
	END
)
WHERE id = $1
RETURNING id;
`

const deleteJobsStatement = `DELETE FROM jobs WHERE project_id = $1 AND id = ANY ($2);`

const cancelJobsStatement = `
UPDATE kenza.jobs
SET status = 'cancelled'
WHERE id = ANY ($2)
RETURNING id;
`
