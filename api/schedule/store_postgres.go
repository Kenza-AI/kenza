package schedule

import (
	"database/sql"

	"github.com/lib/pq"
)

// Postgres is the postgres schedule `Store` implementation.
type Postgres struct {
	DB *sql.DB
}

// GetSchedulesForProject returns a project's schedules.
func (store *Postgres) GetSchedulesForProject(accountID, projectID int64) ([]Schedule, error) {
	const getSchedulesForProjectStatement = `SELECT schedules.id, account_id, project_id, schedules.title, cron_expression, repository, schedules.description, ref
						 FROM schedules
						 INNER JOIN projects ON projects.id = schedules.project_id 
						 WHERE account_id = $1 AND project_id = $2`

	rows, err := store.DB.Query(getSchedulesForProjectStatement, accountID, projectID)
	if err != nil {
		return []Schedule{}, err
	}

	return sqlRowsToSchedules(rows)
}

// Create a schedule.
func (store *Postgres) Create(schedule Schedule) (scheduleID int64, err error) {
	const createScheduleStatement = `INSERT INTO kenza.schedules (project_id, title, ref, cron_expression, description)
					 VALUES ($1, $2, $3, $4, $5)
					 RETURNING id`

	err = store.DB.QueryRow(createScheduleStatement, schedule.ProjectID, schedule.Title, schedule.Ref, schedule.Cron, schedule.Description).Scan(&scheduleID)
	return scheduleID, err
}

// Update a schedule.
func (store *Postgres) Update(schedule Schedule) (scheduleID int64, err error) {
	const updateScheduleStatement = `UPDATE schedules
					 SET description = $3, ref = $4, cron_expression = $5
					 WHERE title = $1 and project_id = $2
					 RETURNING id`
	err = store.DB.QueryRow(updateScheduleStatement, schedule.Title, schedule.ProjectID, schedule.Description, schedule.Ref, schedule.Cron).Scan(&scheduleID)
	return scheduleID, err
}

// Delete a schedule.
func (store *Postgres) Delete(accountID, projectID, scheduleID int64) error {
	const deleteScheduleStatement = `DELETE FROM schedules WHERE id = $1 AND project_id = $2 RETURNING id`

	var id int64
	return store.DB.QueryRow(deleteScheduleStatement, scheduleID, projectID).Scan(&id)
}

// GetAll all schedules, excluding the ids in the excluded list.
func (store *Postgres) GetAll(excludedIDs []int64) ([]Schedule, error) {
	const getSchedulesStatement = `SELECT schedules.id, account_id, project_id, schedules.title, cron_expression, repository, schedules.description, ref
				       FROM schedules
				       INNER JOIN kenza.projects ON schedules.project_id = projects.id
				       WHERE schedules.id <> ALL ($1)`

	rows, err := store.DB.Query(getSchedulesStatement, pq.Array(excludedIDs))
	if err != nil {
		return []Schedule{}, err
	}

	return sqlRowsToSchedules(rows)
}

func sqlRowsToSchedules(rows *sql.Rows) ([]Schedule, error) {
	schedules := []Schedule{}
	for rows.Next() {
		var sch Schedule
		if err := rows.Scan(&sch.ID, &sch.AccountID, &sch.ProjectID, &sch.Title, &sch.Cron, &sch.Repository, &sch.Description, &sch.Ref); err != nil {
			return []Schedule{}, err
		}
		schedules = append(schedules, sch)
	}

	return schedules, nil
}
