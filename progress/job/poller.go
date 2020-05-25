package job

import (
	"encoding/json"

	"github.com/kenza-ai/kenza/event"
	"github.com/kenza-ai/kenza/pubsub"
)

// Start subscribes to and handles incoming "job updated" events.
func Start(store Store, sub pubsub.Subscriber, queue string) {
	if err := sub.Subscribe(queue, event.JobUpdatedRoutingKey, 0, func(body []byte, ack func(ok bool, requeue bool)) {
		ok, requeue := false, false
		defer func() {
			i("acknowledging job update, success '%v', requeue: '%v'", ok, requeue)
			ack(ok, requeue)
		}()

		if err := handleUpdate(body, store); err != nil {
			e(err.Error())
			return
		}
		ok = true
	}); err != nil {
		e("subscriber error: ", err)
	}
	i("subscriber stopping")
}

func handleUpdate(body []byte, store Store) error {
	job := event.JobUpdated{}
	if err := json.Unmarshal(body, &job); err != nil {
		return err
	}

	i("handling job update: %+v", job)

	if err := store.UpdateJob(job); err != nil {
		return err
	}
	return nil
}
