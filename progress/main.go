package main

import (
	"flag"
	"io/ioutil"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/kenza-ai/kenza/logutil"
	"github.com/kenza-ai/kenza/progress/job"
	"github.com/kenza-ai/kenza/pubsub"
)

const (
	apiKeyLocation = "/run/secrets/api_key"
)

var (
	version string
)

var (
	// Access to "job updates" queue access
	queue        = flag.String("job_progress_queue", "kenza.queue.job.updates", "Name of queue that handles job update messages")
	exchange     = flag.String("job_progress_exchange", "kenza.exchange.jobs", "Name of exchange name where job update messages are published on")
	rabbitMQPort = flag.Int64("rabbitmq_port", 5672, "RabbitMQ port")
	rabbitMQHost = flag.String("rabbitmq_host", "pubsub", "RabbitMQ host")
	rabbitMQUser = flag.String("rabbitmq_user", "guest", "RabbitMQ user")
	rabbitMQPass = flag.String("rabbitmq_pass", "guest", "RabbitMQ password")
)

var (
	apiKey      string
	jobs        job.Store
	cleanupOnce sync.Once
	sub         *pubsub.RabbitMQ
)

func init() {
	setupLog()
}

func main() {
	defer cleanup()
	flag.Parse()

	setupGracefulShutdown()
	setupAPIKey()
	setupPubSub()
	setupJobStore()

	job.Start(jobs, sub, *queue)
}

func cleanup() {
	cleanupOnce.Do(func() {
		doCleanup()
	})
}

func doCleanup() {
	logutil.Info("cleaning up")
	if sub != nil {
		logutil.Info("closing connection to exchange")
		if err := sub.Close(); err != nil {
			logutil.Error(err.Error())
		}
	}
}

func setupLog() {
	logutil.Init(os.Stderr, os.Stderr, "progress", version)
}

func setupGracefulShutdown() {
	c := make(chan os.Signal, 1) // Catch various termination signals on channel c
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for sig := range c {
			logutil.Info("received signal: %s", sig.String())
			cleanup()
		}
	}()
}

func setupAPIKey() {
	var err error
	apiKeyData, err := ioutil.ReadFile(apiKeyLocation)
	if err != nil {
		panic(err)
	}
	apiKey = string(apiKeyData)
}

func setupPubSub() {
	var err error
	sub, err = pubsub.NewRabbitMQ(*exchange, *rabbitMQUser, *rabbitMQPass, *rabbitMQHost, *rabbitMQPort, 10)
	if err != nil {
		panic(err)
	}
}

func setupJobStore() {
	var err error
	jobs, err = job.NewJobsClient(apiKey, "v1", version)
	if err != nil {
		panic(err)
	}
}
