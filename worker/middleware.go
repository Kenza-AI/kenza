package main

import "github.com/kenza-ai/kenza/worker/job"

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
