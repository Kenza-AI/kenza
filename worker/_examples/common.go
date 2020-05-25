package examples

import (
	"encoding/json"
	"log"

	"github.com/Kenza-AI/worker/job"
	"github.com/streadway/amqp"
)

// DemoSendMessage — send a job request
func DemoSendMessage() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"kenza-jobs", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	repo := "https://github.com/ilazakis/dummy.git"
	branch := "refs/heads/master"
	commitID := ""
	jobRequest := job.Request{RepoURL: repo, JobID: "1", Branch: branch, CommitID: commitID, ProjectID: "1"}
	body, err := json.Marshal(jobRequest)
	if err != nil {
		log.Print(err)
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
