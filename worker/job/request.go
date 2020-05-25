package job

import (
	"encoding/json"
	"sync"

	"github.com/kenza-ai/kenza/event"
	"github.com/kenza-ai/kenza/pubsub"
)

type Request struct {
	exchange string
	queue    string
	event.JobQueued
	*event.JobUpdated
	pub pubsub.Publisher
	sub pubsub.Subscriber
	*event.Schedules
}

func New(exchange, queue string, pub pubsub.Publisher, sub pubsub.Subscriber) *Request {
	return &Request{exchange: exchange, queue: queue, pub: pub, sub: sub, JobUpdated: &event.JobUpdated{}}
}

func (r *Request) Poll() error {

	i("binding queue %s to exchange %s with routing key %s", r.queue, r.exchange, event.JobQueuedRoutingKey)

	var handleOnce sync.Once
	jobArrival := make(chan event.JobQueued, 1)

	var unmarshallingError error
	go r.sub.Subscribe(r.queue, event.JobQueuedRoutingKey, 1, func(body []byte, ack func(ok bool, requeue bool)) {
		handleOnce.Do(func() {
			defer func() {
				i("closing polling exchange connection")
				r.sub.Close()
			}()
			ack(true, false)

			if err := json.Unmarshal(body, &r.JobQueued); err != nil {
				r.Fail(err)
				unmarshallingError = err
			}

			r.JobUpdated.JobID = r.JobQueued.JobID
			r.JobUpdated.CommitID = r.JobQueued.CommitID
			r.JobUpdated.ProjectID = r.JobQueued.ProjectID
			r.JobUpdated.AccountID = r.JobQueued.AccountID

			i("received job: %+v", r.JobQueued)
			jobArrival <- r.JobQueued
		})
	})

	<-jobArrival
	return unmarshallingError
}

func (r *Request) Notify() error {
	i("Publishing job update %+v", r.JobUpdated)
	return r.pub.Publish(r.JobUpdated, event.JobUpdatedRoutingKey)
}

func (r *Request) NotifySchedules() error {
	if r.Schedules == nil {
		r.Schedules = &event.Schedules{}
	}
	i("Publishing schedule entries for job %d %+v", r.JobQueued.JobID, r.Schedules)
	return r.pub.Publish(event.SchedulesReceived{
		AccountID: r.AccountID,
		ProjectID: r.ProjectID,
		Schedules: r.Schedules,
	}, event.SchedulesReceivedRoutingKey)
}

func (r *Request) Fail(err error) error {
	e("failing job: %s", err)
	r.Status = "failed"
	return r.Notify()
}
