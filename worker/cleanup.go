package main

import (
	"os"
)

func cleanup() {
	if r := recover(); r != nil {
		i("recovering from panic:", r)
	}

	if logfile != nil {
		i("flushing log file buffers before shutdown")
		if err := logfile.Sync(); err != nil {
			e(err.Error())
		}
		logfile.Close()
	}

	if pub != nil {
		i("closing publishing exchange connection")
		if err := pub.Close(); err != nil {
			e(err.Error())
		}
	}

	if sub != nil {
		i("closing polling exchange connection if still alive")
		if err := sub.Close(); err != nil {
			e(err.Error())
		}
	}

	shutdown()
}

// TODO(ilazakis): exit code 0 if job was successfully completed
func shutdown() {
	i("shutting down")
	os.Exit(1)
}
