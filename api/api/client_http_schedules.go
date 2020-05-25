package api

import (
	"fmt"
	"strconv"

	"github.com/kenza-ai/kenza/api/schedule"
)

const schedulesPath = "/accounts/%d/projects/%d/schedules"

// ScheduleCreate creates a schedule to run a job based on a cron expression.
func (c *HTTP) ScheduleCreate(sch schedule.Schedule) (scheduleID int64, err error) {
	path := c.version + fmt.Sprintf(schedulesPath, sch.AccountID, sch.ProjectID)

	i("requesting schedule creation with details %+v", sch)
	req, err := c.newRequest("POST", path, sch)
	if err != nil {
		return -1, err
	}

	var createScheduleResponseBody struct {
		ID int64
	}
	response, err := c.do(req, &createScheduleResponseBody)
	if err != nil {
		i("create schedule response: %+v", response)
		return -1, err
	}

	if response.StatusCode >= 400 {
		err = fmt.Errorf("create schedule error: received status code %d", response.StatusCode)
	}

	i("create schedule response: %+v", response)
	return createScheduleResponseBody.ID, err
}

// ScheduleUpdate updates a schedule.
func (c *HTTP) ScheduleUpdate(sch schedule.Schedule) error {
	path := c.version + fmt.Sprintf(schedulesPath, sch.AccountID, sch.ProjectID)

	i("requesting schedule update with details %+v", sch)
	req, err := c.newRequest("PUT", path, sch)
	if err != nil {
		return err
	}

	var updateScheduleResponseBody struct {
		ID int64
	}
	response, err := c.do(req, &updateScheduleResponseBody)
	if err != nil {
		i("update schedule response: %+v", response)
		return err
	}

	if response.StatusCode >= 400 {
		err = fmt.Errorf("update schedule error: received status code %d", response.StatusCode)
	}

	i("update schedule response: %+v", response)
	return err
}

// ScheduleDelete deletes a schedule.
func (c *HTTP) ScheduleDelete(accountID, projectID, scheduleID int64) error {
	path := c.version + fmt.Sprintf(schedulesPath+"/%d", accountID, projectID, scheduleID)

	i("requesting schedule deletion with details %+v", scheduleID)
	req, err := c.newRequest("DELETE", path, scheduleID)
	if err != nil {
		return err
	}

	response, err := c.do(req, nil)
	if err != nil {
		i("delete schedule response: %+v", response)
		return err
	}

	if response.StatusCode >= 400 {
		err = fmt.Errorf("delete schedule error: received status code %d", response.StatusCode)
	}

	i("delete schedule response: %+v", response)
	return err
}

// SchedulesForProject returns a project's schedules.
func (c *HTTP) SchedulesForProject(accountID, projectID int64) ([]schedule.Schedule, error) {
	path := c.version + fmt.Sprintf(schedulesPath, accountID, projectID)

	i("requesting schedules for project '%d' in account '%d'", projectID, accountID)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return []schedule.Schedule{}, err
	}

	var schedulesForProjectResponseBody struct {
		Schedules []schedule.Schedule
	}
	response, err := c.do(req, &schedulesForProjectResponseBody)
	if err != nil {
		i("get schedules response: %+v", response)
		return []schedule.Schedule{}, err
	}

	if response.StatusCode >= 400 {
		err = fmt.Errorf("get schedules error: received status code %d", response.StatusCode)
	}

	i("get schedules response: %+v", response)
	return schedulesForProjectResponseBody.Schedules, err
}

// Schedules returns all schedules outside the ones matching the exclusion list of schedule ids.
func (c *HTTP) Schedules(excludedIDs []int64) ([]schedule.Schedule, error) {
	path := c.version + fmt.Sprintf("/schedules")

	i("requesting all schedules excluding ids in: %v", excludedIDs)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return []schedule.Schedule{}, err
	}

	params := req.URL.Query()
	for _, scheduleID := range excludedIDs {
		params.Add("excluded", strconv.FormatInt(scheduleID, 10))
	}
	req.URL.RawQuery = params.Encode()

	var schedulesResponseBody struct {
		Schedules []schedule.Schedule
	}
	response, err := c.do(req, &schedulesResponseBody)
	if err != nil {
		i("get schedules response: %+v", response)
		return []schedule.Schedule{}, err
	}

	if response.StatusCode >= 400 {
		err = fmt.Errorf("get schedules error: received status code %d", response.StatusCode)
	}

	i("get schedules response: %+v", response)
	return schedulesResponseBody.Schedules, err
}
