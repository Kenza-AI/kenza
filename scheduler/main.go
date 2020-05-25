package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/kenza-ai/kenza/api/api"
	"github.com/kenza-ai/kenza/logutil"
	"github.com/kenza-ai/kenza/pubsub"
	"github.com/kenza-ai/kenza/scheduler/job"
	"github.com/robfig/cron/v3"
)

const (
	apiKeyLocation = "/run/secrets/api_key"
)

var (
	version string
)

var (
	// Broker / pubsub
	queue          = flag.String("queued_jobs_queue", "kenza.queue.jobs.triage", "Name of queue where jobs land on to be processed")
	schedulesQueue = flag.String("schedules_queue", "kenza.queue.schedules.triage", "Name of queue where schedules land on to be processed")
	exchange       = flag.String("queued_jobs_exchange", "kenza.exchange.jobs", "Name of exchange name where 'job queued' messages are published on")
	rabbitMQPort   = flag.Int64("rabbitmq_port", 5672, "RabbitMQ port")
	rabbitMQHost   = flag.String("rabbitmq_host", "pubsub", "RabbitMQ host")
	rabbitMQUser   = flag.String("rabbitmq_user", "guest", "RabbitMQ user")
	rabbitMQPass   = flag.String("rabbitmq_pass", "guest", "RabbitMQ password")
)

var (
	cleanupOnce sync.Once
	apiClient   api.Client
	pub         *pubsub.RabbitMQ
	sub         *pubsub.RabbitMQ
)

var (
	i = logutil.Info
	e = logutil.Error
)

func init() {
	setupLog()
	flag.Parse()
}

func main() {
	defer cleanup()
	setupGracefulShutdown()

	setupPubSub()
	setupAPIClient()
	listenAndHandle()
}

func listenAndHandle() {
	err := make(chan error)

	c := cron.New()
	go job.ListenAndHandleOnDemandJobs(apiClient, sub, pub, *queue, *exchange, err)
	go job.ListenAndHandleScheduledJobs(apiClient, pub, c)
	go job.ListenAndHandleSchedules(apiClient, sub, *schedulesQueue, *exchange, err, c)

	<-err

	i("Scheduler shutting down")
}

func setupLog() {
	logutil.Init(os.Stderr, os.Stderr, "scheduler", version)
}

func setupAPIClient() {
	var err error
	apiKeyData, err := ioutil.ReadFile(apiKeyLocation)
	if err != nil {
		panic(err)
	}

	apiClient, err = newAPIClient(string(apiKeyData), "v1", version)
	if err != nil {
		panic(err)
	}
}

func setupPubSub() {
	var err error
	pub, err = pubsub.NewRabbitMQ(*exchange, *rabbitMQUser, *rabbitMQPass, *rabbitMQHost, *rabbitMQPort, 10)
	if err != nil {
		panic(err)
	}

	sub, err = pubsub.NewRabbitMQ(*exchange, *rabbitMQUser, *rabbitMQPass, *rabbitMQHost, *rabbitMQPort, 10)
	if err != nil {
		panic(err)
	}
}

func cleanup() {
	cleanupOnce.Do(func() {
		doCleanup()
	})
}

func doCleanup() {
	i("cleaning up")
	if sub != nil {
		i("closing subscriber connection to exchange")
		if err := sub.Close(); err != nil {
			e(err.Error())
		}
	}

	if pub != nil {
		i("closing publisher connection to exchange")
		if err := pub.Close(); err != nil {
			e(err.Error())
		}
	}
}

func setupGracefulShutdown() {
	c := make(chan os.Signal, 1) // Catch various termination signals on channel c
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for sig := range c {
			i("received signal: %s", sig.String())
			cleanup()
		}
	}()
}

// NewAPIClient wraps a `client.HTTP` API `Client` and forwards requests
// to the Kenza API (as opposed to directly depending on a store / DB).
func newAPIClient(apiKey, apiVersion, serviceVersion string) (*api.HTTP, error) {
	httpClient := &http.Client{
		Timeout: time.Second * 15,
	}
	userAgent := "Scheduler/" + serviceVersion

	return api.NewHTTPClient("http://api:8080", apiVersion, userAgent, apiKey, httpClient)
}
