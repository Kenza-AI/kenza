package main

import (
	. "github.com/Kenza-AI/worker/_examples"
	"github.com/Kenza-AI/worker/job"
)

func main() {
	DemoSendMessage()

	// poller := initialization.NewRabbitMQPoller()
	// job := poller.Start()

	// versionControl := handler.NewVCS(&vcs.Git{})
	// service := handler.NewService(".kenza.yml")

	// notifier := &notifier.Log{}

	// h := chainHandlers(versionControl, service)

	// if len(h) > 0 {
	// 	h[0].Handle(notifier, *job)
	// }
}

func chainHandlers(handlers ...job.ChainHandler) []job.ChainHandler {
	len := len(handlers)
	if len < 1 {
		return []job.ChainHandler{}
	}

	for idx, h := range handlers {
		hasNext := len > idx+1
		if hasNext {
			h.SetNext(handlers[idx+1])
		}
	}
	return handlers
}
