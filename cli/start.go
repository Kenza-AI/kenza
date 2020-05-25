package cli

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const startLongDescription = `
The start command runs the Kenza service(s) passed as args. 

Passing no arguments starts the whole stack. That's the most common scenario. 

Running individual services is typically used to diagnose unexpected behavior or during dev.

Valid services: [web, api, progress, worker, db, pubsub].
`

const startExamples = `
To start only the database and api services run:
kenza start api db 
`

var (
	apiKey              string
	machineName         string
	githubWebhookSecret string
)

func init() {
	parseStartCmdFlags()
}

var startCmd = &cobra.Command{
	Use:     "start",
	Short:   "Deploys or updates the Kenza Docker stack",
	Long:    startLongDescription,
	Example: startExamples,
	Args:    cobra.NoArgs,

	Run: func(cmd *cobra.Command, args []string) {
		ackMessage := "Deploying Kenza locally"
		if machineName != "" {
			ackMessage = fmt.Sprintf("Deploying Kenza on AWS EC2 for installation named '%s'", machineName)
		}
		fmt.Println(ackMessage)

		initCluster()

		// Create directories kenza expects
		if err := prepareStartCommandDirs(); err != nil {
			panic(err)
		}

		// Create docker-compose.yml, .env and secrets files
		if err := prepareStartCommandFiles(); err != nil {
			panic(err)
		}

		// Pull all Docker images first. `docker stack deploy` pulls images too
		// but there's currently no way to get downloading progress feedback.
		// `docker pull` does send progress feedback to stdout instead, so we pull first.
		if err := pullImages(); err != nil {
			panic(err)
		}

		// Start / deploy services
		if output, err := executeNoPrint(startCommand(args)); err != nil {
			fmt.Println(output)
			panic(err)
		}

		// Open browser to Kenza web app
		hostname, err := hostname()
		if err != nil {
			panic(err)
		}

		fmt.Printf("\n%s Host IP: %s\n", machineName, hostname)
		if err := open(hostname); err != nil {
			fmt.Println(err)
		}
	},
}

func initCluster() {
	initClusterCommand := []string{"docker swarm init"}

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	if err := executeWithErrOut(initClusterCommand, stderr, stdout); err != nil {
		trimmedStderr := strings.Trim(stderr.String(), "\n")
		if trimmedStderr != errSwarmExists {
			panic(stderr.String())
		}
		fmt.Print("Already part of a cluster (docker swarm). Will create and update services as needed.\n\n")
		return
	}
	fmt.Print(stdout.String())
}

func pullImages() error {
	imageRegistry := "aikenza"
	for _, service := range kenzaServices {

		// Postgres and rabbitmq do not live in Kenza's registry.
		// Pulling from their public registries instead.
		dockerPullCommand := ""
		switch service {
		case "db":
			dockerPullCommand = "docker pull postgres"
		case "pubsub":
			dockerPullCommand = "docker pull rabbitmq:3-management"
		default:
			dockerPullCommand = fmt.Sprintf("docker pull %s/%s:%s", imageRegistry, service, Version)
		}

		fmt.Printf("Pulling '%s' service image (%s)\n", service, dockerPullCommand)
		if err := execute([]string{dockerPullCommand}); err != nil {
			return err
		}
	}
	return nil
}

func startCommand(args []string) []string {

	// With the directories and compose files in place under ./kenza
	// we switch to that directory to get an updated compose file with
	// all the env variables substituted (via docker-compose config)
	err := os.Chdir("kenza")
	if err != nil {
		panic(err)
	}

	// Set version env var for compose file to pick up the corresponding image tags when building
	os.Setenv("KENZA_VERSION", Version)

	if inCloudMode() {
		// Set API host (needed by the api service to point source control webhooks to Kenza)
		hostname, err := hostname()
		if err != nil {
			panic(err)
		}
		os.Setenv("KENZA_API_HOST", hostname)
	}

	cmd := exec.Command("sh", "-c", "docker-compose config")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	if err := writeFile("docker-compose.yml", string(out)); err != nil {
		panic(err)
	}

	err = os.Chdir("..")
	if err != nil {
		panic(err)
	}

	command := []string{"docker stack deploy --with-registry-auth -c kenza/docker-compose.yml kenza"}

	fmt.Println("Starting Kenza...")
	return command
}

func prepareStartCommandDirs() error {
	// directories expected to be present by docker swarm. They should match the corresponding ones in .env
	kenzaInstallationDirectories := []string{
		"./kenza/data/jobs/logs/",
		"./kenza/data/postgres/storage/",
		"./kenza/data/rabbitmq/storage/",
	}

	for _, dir := range kenzaInstallationDirectories {
		if err := os.MkdirAll(filepath.Dir(dir), os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func prepareStartCommandFiles() error {
	for file, content := range dockerComposeFiles() {
		if err := writeFile(file, content); err != nil {
			return err
		}
	}

	if inCloudMode() {
		fmt.Println("Uploading Kenza files to remote cluster")
		return copyKenzaFilesToMachine(machineName)
	}
	return nil
}

func copyKenzaFilesToMachine(machineName string) error {
	return execute([]string{fmt.Sprintf("docker-machine scp -r kenza/ %s:kenza/", machineName)})
}

const (
	// errSwarmExists is returned on initialize or join request for a cluster that has already been activated
	errSwarmExists string = "Error response from daemon: This node is already part of a swarm. Use \"docker swarm leave\" to leave this swarm and join another one."
)

func hostname() (hostname string, err error) {
	if !inCloudMode() {
		return "localhost", nil
	}

	// Get machine/manager IP address / hostname
	output, err := executeNoPrint([]string{fmt.Sprintf("docker-machine ip %s", machineName)})
	if err != nil {
		fmt.Println(output)
		return "", err
	}
	return strings.TrimSpace(output), nil
}

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}

	// TODO(ilazakis): https support
	args = append(args, "http://"+url)

	// TODO(ilazakis): poll web container until healthy instead of guestimating when it will be up
	fmt.Println("\nWaiting for web container to start...")
	time.Sleep(time.Second * 15)

	fmt.Println("\nOpening default browser at URL: ", args)
	return exec.Command(cmd, args...).Start()
}

func dockerComposeFiles() map[string]string {
	envFile, composeFile := env, compose
	if inCloudMode() {
		envFile = envCloud
		composeFile = composeCloud
	}
	return map[string]string{
		// Compose and environment config
		"kenza/.env":               envFile,
		"kenza/docker-compose.yml": composeFile,

		// Secrets
		"kenza/api_key.secret":               apiKey,
		"kenza/github_webhook_secret.secret": githubWebhookSecret,
	}
}

func inCloudMode() bool {
	return machineName != ""
}

func parseStartCmdFlags() {
	startCmd.Flags().StringVarP(&machineName, "name", "n", "", "Name for the Kenza installation to start / deploy.")
	startCmd.Flags().StringVarP(&apiKey, "apikey", "", "default_api_key_change_when_deploying_to_cloud", "Service-to-Service API key,presented by Kenza services to the API service.")
	startCmd.Flags().StringVarP(&githubWebhookSecret, "github-secret", "", "", "GitHub secret used to secure and authenticate webhook requests.")
}
