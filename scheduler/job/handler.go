package job

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/kenza-ai/kenza/api/api"
	"github.com/kenza-ai/kenza/api/schedule"
	"github.com/kenza-ai/kenza/event"
	"github.com/kenza-ai/kenza/pubsub"
	"github.com/robfig/cron/v3"
)

var (
	scheduledMu      sync.RWMutex
	scheduleCronJobs = map[int64]cron.EntryID{}
)

// ListenAndHandleOnDemandJobs handles on-demand jobs (jobs arriving via a VCS webhook or from the UI directly).
func ListenAndHandleOnDemandJobs(apiClient api.Client, pub, sub *pubsub.RabbitMQ, queue, exchange string, done chan<- error) {
	i("Binding queue %s to exchange %s on routing key %s", queue, exchange, event.JobArrivedRoutingKey)

	err := sub.Subscribe(queue, event.JobArrivedRoutingKey, 0, func(body []byte, ack func(ok bool, requeue bool)) {
		var ok bool
		defer func() {
			i("Enqueued '%v', requeued: '%v'", ok, false)
			ack(ok, false)
		}()

		arrival := event.JobArrived{}
		if err := json.Unmarshal(body, &arrival); err != nil {
			e("Incoming job request error %s", err)
			return
		}
		i("Received job request %+v", arrival)

		jobID, err := apiClient.JobCreate(arrival.AccountID, arrival.ProjectID, arrival.Submitter, arrival.DeliveryID, arrival.CommitID)
		if err != nil {
			e(err.Error())
			return
		}
		i(`Created job "%d" in project "%d" in account "%d"`, jobID, arrival.ProjectID, arrival.AccountID)

		queuedEvt := event.JobQueued{JobID: jobID, JobInfo: arrival}
		if err := pub.Publish(queuedEvt, event.JobQueuedRoutingKey); err != nil {
			e("Failed to enqueue job '%s'", err)
			return
		}
		ok = true
		i("Enqueued job for processing %+v", queuedEvt)
	})

	if err != nil {
		e("Jobs subscriber error: ", err)
	}

	i("Job arrivals subscriber stopping")
	done <- err
}

// ListenAndHandleScheduledJobs handles scheduled jobs (jobs starting from a cron job).
// It regularly polls for schedules to detect which ones have not yet been submitted and
// submits cron jobs for those schedules that have not yet been setup to run.
func ListenAndHandleScheduledJobs(apiClient api.Client, pub *pubsub.RabbitMQ, c *cron.Cron) {
	c.Start()

	ticker := time.NewTicker(time.Second * 60)
	defer ticker.Stop()

	// For each not already submitted schedule, start a cron job that
	// will create and submit a Kenza job for processing to the workers.
	for range ticker.C {

		scheduledMu.RLock()
		submittedSchedules := make([]int64, len(scheduleCronJobs))
		for cronJobEntryID := range scheduleCronJobs {
			submittedSchedules = append(submittedSchedules, cronJobEntryID)
		}
		scheduledMu.RUnlock()

		schedules, err := apiClient.Schedules(submittedSchedules)
		if err != nil {
			e("error retrieving schedules: %s", err.Error())
			continue
		}

		count := 0
		for _, sch := range schedules {
			i("Retrieved schedule %d", sch.ID)

			cronEntryID, err := c.AddFunc(sch.Cron, func() { setupScheduledJob(apiClient, pub, sch) })
			if err != nil {
				e("Cron error: %s", err)
				continue
			}

			count++
			scheduledMu.Lock()
			scheduleCronJobs[sch.ID] = cronEntryID
			scheduledMu.Unlock()
			i("Submitted schedule %d with cron id %d", sch.ID, cronEntryID)
		}
		i("Added cron for %d new schedule(s)", count)
	}
}

func setupScheduledJob(apiClient api.Client, pub *pubsub.RabbitMQ, sch schedule.Schedule) {
	i("Running cron for schedule: %+v", sch)

	// Create job
	jobID, err := apiClient.JobCreate(sch.AccountID, sch.ProjectID, sch.Title, "DELIVERY-ID TODO", "")
	if err != nil {
		e(err.Error())
		return
	}
	i(`Created job "%d" in project "%d" in account "%d"`, jobID, sch.ProjectID, sch.AccountID)

	// Submit job to worker queue
	evt := event.JobQueued{
		JobID: jobID,
		JobInfo: event.JobInfo{
			AccountID: sch.AccountID,
			ProjectID: sch.ProjectID,
			CloneURL:  sch.Repository,
			Ref:       sch.Ref,
			CommitID:  "", // (HEAD or the commit a tag is referring to)
		},
	}

	if err := pub.Publish(evt, event.JobQueuedRoutingKey); err != nil {
		e("Failed to enqueue job '%s'", err)
		return // TODO(ilazakis): attempt to save job or drop in a DLX
	}
	i("Enqueued job %d", jobID)
}

// ListenAndHandleSchedules handles changes to schedules arriving from job workers (via changes to .kenza.yml)
func ListenAndHandleSchedules(apiClient api.Client, sub *pubsub.RabbitMQ, schedulesQueue, exchange string, done chan<- error, c *cron.Cron) {
	i("Binding queue %s to exchange %s on routing key %s", schedulesQueue, exchange, event.SchedulesReceivedRoutingKey)

	err := sub.Subscribe(schedulesQueue, event.SchedulesReceivedRoutingKey, 0, func(body []byte, ack func(ok bool, requeue bool)) {
		var ok bool
		defer func() {
			i("Enqueued '%v', requeued: '%v'", ok, false)
			ack(ok, false)
		}()

		incomingSchedulesEvent := &event.SchedulesReceived{}
		if err := json.Unmarshal(body, incomingSchedulesEvent); err != nil {
			e("Incoming schedules request error %s", err)
			return
		}
		i("Received schedules event for account %d project %d, number of schedules %d",
			incomingSchedulesEvent.AccountID,
			incomingSchedulesEvent.ProjectID,
			len(*incomingSchedulesEvent.Schedules))
		incomingSchedules := incomingSchedulesEvent.Schedules

		existingSchedules, err := apiClient.SchedulesForProject(incomingSchedulesEvent.AccountID, incomingSchedulesEvent.ProjectID)
		if err != nil {
			e("error retrieving schedules %s", err.Error())
			return
		}
		i("retrieved schedules %+v", existingSchedules)

		handleSchedules(apiClient, existingSchedules, *incomingSchedules, incomingSchedulesEvent.AccountID, incomingSchedulesEvent.ProjectID, c)
		ok = true
	})

	if err != nil {
		e("Schedules subscriber error: ", err)
	}

	i("Schedules subscriber stopping")

	done <- err
}

// 1. Deletes schedules not in `incoming` schedules
// 2. Updates schedules present in both `existing` and `incoming` schedules
// 3. Creates schedules only present in `incoming` schedules
func handleSchedules(apiClient api.Client, existingSchedules []schedule.Schedule, incomingSchedules event.Schedules, accountID, projectID int64, c *cron.Cron) {
	existingSchedulesTitles := map[string]struct{}{}
	for _, schedule := range existingSchedules {
		existingSchedulesTitles[schedule.Title] = struct{}{}

		_, ok := incomingSchedules[schedule.Title]
		if !ok {
			err := apiClient.ScheduleDelete(accountID, projectID, schedule.ID)
			if err != nil {
				e(err.Error())
				continue
			}
			i("Schedule deleted %d", schedule.ID)

			scheduledMu.RLock()
			cronForRemoval, ok := scheduleCronJobs[schedule.ID]
			if !ok {
				continue
			}
			scheduledMu.Unlock()
			c.Remove(cronForRemoval)
		}
	}

	for incomingScheduleTitle, scheduleEntry := range incomingSchedules {
		_, alreadyExisted := existingSchedulesTitles[incomingScheduleTitle]

		ref := "refs/heads/" + scheduleEntry.Branch
		if scheduleEntry.Tag != "" {
			ref = "refs/tags/" + scheduleEntry.Tag
		}

		var err error
		var schID int64
		if alreadyExisted {
			err = apiClient.ScheduleUpdate(schedule.Schedule{
				AccountID:   accountID,
				ProjectID:   projectID,
				Ref:         ref,
				Cron:        scheduleEntry.When,
				Title:       incomingScheduleTitle,
				Description: scheduleEntry.Description,
			})
			if err != nil {
				e("error updating schedule: %s", err.Error())
				continue
			}
			scheduledMu.Lock()
			cronForRemoval, ok := scheduleCronJobs[schID]
			if !ok {
				scheduledMu.Unlock()
				continue
			}
			c.Remove(cronForRemoval)
			delete(scheduleCronJobs, schID)
			scheduledMu.Unlock()
		} else {
			schID, err = apiClient.ScheduleCreate(schedule.Schedule{
				AccountID:   accountID,
				ProjectID:   projectID,
				Ref:         ref,
				Cron:        scheduleEntry.When,
				Title:       incomingScheduleTitle,
				Description: scheduleEntry.Description,
			})
			if err != nil {
				e("error creating schedule: %s", err.Error())
				continue
			}
		}

		if alreadyExisted {
			i("Updated schedule %d and removed pending scheduled jobs for updates to be taken into account next time cron runs", schID)
		} else {
			i("Created schedule %d", schID)
		}
	}
}
