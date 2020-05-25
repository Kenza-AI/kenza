package main

import (
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/kenza-ai/kenza/logutil"
	"github.com/kenza-ai/kenza/pubsub"
	"github.com/kenza-ai/kenza/worker/job"
	"github.com/kenza-ai/kenza/worker/job/handler"
	"github.com/kenza-ai/kenza/worker/secrets"
	"github.com/kenza-ai/kenza/worker/vcs"
)

var (
	version string
)

const (
	apiKeyLocation = "/run/secrets/api_key"
)

var (
	apiKey  string
	logfile *os.File
	j       *job.Request
	sec     secrets.Store
	pub     *pubsub.RabbitMQ
	sub     *pubsub.RabbitMQ
)

var (
	// Access to jobs "queued" for processing
	queue        = flag.String("queued_jobs_queue", "kenza.queue.jobs.ready", "Name of queue where jobs land on to be processed")
	exchange     = flag.String("queued_jobs_exchange", "kenza.exchange.jobs", "Name of exchange name where job update messages are published on")
	rabbitMQPort = flag.Int64("rabbitmq_port", 5672, "RabbitMQ port")
	rabbitMQHost = flag.String("rabbitmq_host", "pubsub", "RabbitMQ host")
	rabbitMQUser = flag.String("rabbitmq_user", "guest", "RabbitMQ user")
	rabbitMQPass = flag.String("rabbitmq_pass", "guest", "RabbitMQ password")

	// Job artifacts locations
	workDir    = flag.String("work_dir", os.TempDir(), "Where the repo will be cloned")
	logfileDir = flag.String("logfile_dir", "", "Path where job logs are written")
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
	setupAPIClient()
	listenAndHandle(setupHandlers())
}

func setupLog() {
	logutil.Init(os.Stderr, os.Stderr, "worker", version)
}

func setupLogfile(jobID, projectID int64) {
	var err error
	logfilePath := filepath.Join(*logfileDir, strconv.FormatInt(projectID, 10)+"-"+strconv.FormatInt(jobID, 10)+".log")

	logfile, err = os.OpenFile(logfilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		e("could not create log file at '%s' — logging to stdout only. %s", *logfileDir, err)
	} else {
		i("created log file at '%s'", logfilePath)
		mw := io.MultiWriter(os.Stdout, logfile)
		logutil.SetOutputErr(mw)
		logutil.SetOutputInfo(mw)
	}
}

func listenAndHandle(h []job.ChainHandler) {
	j = job.New(*exchange, *queue, pub, sub)
	if err := j.Poll(); err != nil {
		panic(err)
	}

	i("Creating log file for job %+v", j.JobQueued.JobID)
	setupLogfile(j.JobQueued.JobID, j.JobQueued.ProjectID)

	if len(h) > 0 {
		h[0].Handle(j)
	}
}

func setupAPIClient() {
	var err error
	sec, err = secrets.NewSecretsClient(apiKey, "v1", version)
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

func setupHandlers() []job.ChainHandler {
	arrival := handler.NewArrival(*workDir)
	versionControl := handler.NewVCS(&vcs.Git{}, sec)
	service := handler.NewService()
	return chainHandlers(arrival, versionControl, service)
}

func setupGracefulShutdown() {
	c := make(chan os.Signal, 1) // Catch various termination signals on channel c
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for sig := range c {
			i("received signal: %s", sig.String())
			if j != nil {
				j.Fail(errors.New(sig.String()))
			}
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
