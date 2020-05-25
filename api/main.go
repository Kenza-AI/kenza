package main

import (
	"context"
	"database/sql"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/kenza-ai/kenza/api/account"
	"github.com/kenza-ai/kenza/api/httputil"
	"github.com/kenza-ai/kenza/api/job"
	"github.com/kenza-ai/kenza/api/project"
	"github.com/kenza-ai/kenza/api/sagemaker"
	"github.com/kenza-ai/kenza/api/schedule"
	"github.com/kenza-ai/kenza/db"
	"github.com/kenza-ai/kenza/logutil"
	"github.com/kenza-ai/kenza/pubsub"
	"github.com/rs/cors"

	"github.com/google/go-github/v29/github"
)

var (
	version string
)

const (
	apiKeyLocation              = "/run/secrets/api_key"
	githubWebHookSecretLocation = "/run/secrets/github_webhook_secret"
)

var (
	// Database access
	dbPort = flag.Int64("db_port", 5432, "DB port")
	dbHost = flag.String("db_host", "db", "DB host")
	dbName = flag.String("db_name", "kenza", "DB name")
	dbUser = flag.String("db_user", "kenza", "DB user")
	dbPass = flag.String("db_pass", "kenza", "DB password")

	// Queued jobs access
	exchange     = flag.String("queued_jobs_exchange", "kenza.exchange.jobs", "Name of exchange name where 'job queued' messages are published on")
	rabbitMQPort = flag.Int64("rabbitmq_port", 5672, "RabbitMQ port")
	rabbitMQHost = flag.String("rabbitmq_host", "pubsub", "RabbitMQ host")
	rabbitMQUser = flag.String("rabbitmq_user", "guest", "RabbitMQ user")
	rabbitMQPass = flag.String("rabbitmq_pass", "guest", "RabbitMQ password")

	// AUTH
	apiKey              = "" // read from mounted secrets directory
	githubWebHookSecret = "" // read from mounted secrets directory
	jwtSigningKey       = flag.String("jwt_signing_key", "", "JWT singing key. Resetting logs everyone out.")

	// AWS
	awsProfile         = flag.String("aws_profile", "", "AWS profile used to read SageMaker job details.")
	awsConfigFile      = flag.String("aws_config", "", "AWS configuration file location.")
	awsCredentialsFile = flag.String("aws_credentials", "", "AWS credentials file location.")

	// Job logs location
	logfileDir = flag.String("logfile_dir", "", "Path to job logs")

	// Listening address for source control webhooks to point at
	port = flag.String("api_port", ":8080", "Listening port for source control webhooks to point at.")
	host = flag.String("api_host", "localhost", "Listening host for source control webhooks to point at.")
)

var (
	database *sql.DB
	sm       *sagemaker.Client

	jobs      job.Store
	projects  project.Store
	accounts  account.Store
	schedules schedule.Store

	pub *pubsub.RabbitMQ

	srv         http.Server
	inShutdown  int32
	doneCleanup = make(chan bool, 1)
)

func init() {
	setupLog()
}

func main() {
	defer func() {
		cleanup()
	}()
	flag.Parse()

	setupDB()
	setupPubSub()
	setupSageMaker()
	setupServer()
	startServer()

	<-doneCleanup
	logutil.Info("server stopped")
}

func cleanup() {
	if shuttingDown() {
		return
	}
	atomic.StoreInt32(&inShutdown, 1)
	defer close(doneCleanup)

	logutil.Info("cleaning up")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	logutil.Info("shutting server down")
	srv.SetKeepAlivesEnabled(false)
	if err := srv.Shutdown(ctx); err != nil {
		logutil.Error("server shutdown error: %v", err)
	}

	if pub != nil {
		logutil.Info("closing connection to exchange")
		if err := pub.Close(); err != nil {
			logutil.Error(err.Error())
		}
	}

	if database != nil {
		logutil.Info("closing connection to database")
		if err := database.Close(); err != nil {
			logutil.Error(err.Error())
		}
	}

	logutil.Info("finished cleaning up")
}

func setupServer() {
	srv = http.Server{
		Addr:              *port,
		Handler:           routes(),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       20 * time.Second,
		WriteTimeout:      45 * time.Second,
	}
}

func setupAPIKey() {
	apiKeyData, err := ioutil.ReadFile(apiKeyLocation)
	if err != nil {
		panic(err)
	}
	apiKey = string(apiKeyData)

	// Local mode does not use GitHub webhooks
	if *host == "localhost" {
		return
	}

	githubSecretData, err := ioutil.ReadFile(githubWebHookSecretLocation)
	if err != nil {
		panic(err)
	}
	githubWebHookSecret = string(githubSecretData)
}

func setupSageMaker() {
	var err error
	sm, err = sagemaker.NewClient(*awsProfile, "us-east-1", "kenza-api", []string{*awsConfigFile, *awsCredentialsFile})
	if err != nil {
		logutil.Error(err.Error())
	}
}

func setupDB() {
	var err error
	database, err = db.New(*dbUser, *dbPass, *dbHost, *dbName, *dbPort, true)
	if err != nil {
		panic(err)
	}
	jobs = &job.Postgres{DB: database}
	projects = &project.Postgres{DB: database}
	accounts = &account.Postgres{DB: database}
	schedules = &schedule.Postgres{DB: database}
}

func setupPubSub() {
	var err error
	pub, err = pubsub.NewRabbitMQ(*exchange, *rabbitMQUser, *rabbitMQPass, *rabbitMQHost, *rabbitMQPort, 10)
	if err != nil {
		panic(err)
	}
}

func setupLog() {
	logutil.Init(os.Stderr, os.Stderr, "api", version)
}

func routes() http.Handler {
	setupAPIKey()
	router := httprouter.New()
	httpClient := &http.Client{Timeout: time.Second * 30}

	log, auth := httputil.Log, account.Authorize
	GET, POST, PUT, DELETE := http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete

	webhookValidator := github.ValidatePayload

	c := cors.New(cors.Options{
		AllowedMethods: []string{GET, POST, PUT, DELETE},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}).Handler(router)

	// v1 accounts / auth
	router.Handler(POST, "/v1/users", log(account.SignUp(accounts)))
	router.Handler(POST, "/v1/tokens", log(account.SignIn(accounts, *jwtSigningKey)))

	// v1 jobs
	router.Handler(POST, "/v1/jobs/submissions", log(job.SubmitJob(jobs, projects, pub, webhookValidator, githubWebHookSecret))) // endpoint used by VCS webhooks e.g. GitHub
	router.Handler(POST, "/v1/accounts/:accountID/projects/:projectID/jobs/submissions", log(auth(job.SubmitJob(jobs, projects, pub, webhookValidator, githubWebHookSecret), *jwtSigningKey, apiKey)))

	router.Handler(POST, "/v1/accounts/:accountID/projects/:projectID/jobs", log(auth(job.CreateJob(jobs, projects), *jwtSigningKey, apiKey)))
	router.Handler(DELETE, "/v1/accounts/:accountID/projects/:projectID/jobs", log(auth(job.Delete(jobs), *jwtSigningKey, apiKey)))

	router.Handler(GET, "/v1/accounts/:accountID/projects/:projectID/jobs", log(auth(job.GetAll(jobs), *jwtSigningKey, apiKey)))
	router.Handler(GET, "/v1/accounts/:accountID/projects/:projectID/jobs/:jobID", log(auth(job.Get(jobs, sm), *jwtSigningKey, apiKey)))
	router.Handler(GET, "/v1/accounts/:accountID/projects/:projectID/jobs/:jobID/logs", log(auth(job.GetLog(*logfileDir), *jwtSigningKey, apiKey)))

	router.Handler(PUT, "/v1/accounts/:accountID/projects/:projectID/jobs/:jobID", log(auth(job.Put(jobs), *jwtSigningKey, apiKey)))
	router.Handler(POST, "/v1/accounts/:accountID/projects/:projectID/jobs/cancellations", log(auth(job.Cancel(jobs), *jwtSigningKey, apiKey)))

	// v1 projects
	router.Handler(DELETE, "/v1/accounts/:accountID/projects/:projectID", log(auth(project.Delete(projects), *jwtSigningKey, apiKey)))
	router.Handler(GET, "/v1/accounts/:accountID/projects", log(auth(project.GetAll(projects), *jwtSigningKey, apiKey)))
	router.Handler(POST, "/v1/accounts/:accountID/projects", log(auth(project.Create(projects, httpClient, *host+*port, githubWebHookSecret), *jwtSigningKey, apiKey)))
	router.Handler(GET, "/v1/accounts/:accountID/projects/:projectID/secrets", log(auth(project.GetAccessToken(projects), *jwtSigningKey, apiKey)))

	// v1 schedules
	router.Handler(PUT, "/v1/accounts/:accountID/projects/:projectID/schedules", log(auth(schedule.Update(schedules), *jwtSigningKey, apiKey)))
	router.Handler(POST, "/v1/accounts/:accountID/projects/:projectID/schedules", log(auth(schedule.Create(schedules), *jwtSigningKey, apiKey)))
	router.Handler(DELETE, "/v1/accounts/:accountID/projects/:projectID/schedules/:scheduleID", log(auth(schedule.Delete(schedules), *jwtSigningKey, apiKey)))
	router.Handler(GET, "/v1/accounts/:accountID/projects/:projectID/schedules", log(auth(schedule.GetSchedulesForProject(schedules), *jwtSigningKey, apiKey)))
	router.Handler(GET, "/v1/schedules", log(auth(schedule.GetSchedules(schedules), *jwtSigningKey, apiKey)))

	return c
}

func startServer() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGQUIT)

	go handleSignals(c)

	logutil.Info("db, pubsub setup complete, starting listening for connections")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logutil.Error("server error: %v", err)
		cleanup()
	}
}

func handleSignals(c chan os.Signal) {
	sig := <-c
	logutil.Info("received signal: %s", sig.String())
	cleanup()
}

func shuttingDown() bool {
	return atomic.LoadInt32(&inShutdown) == 1
}
