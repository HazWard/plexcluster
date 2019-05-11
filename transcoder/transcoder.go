package transcoder

import (
	"encoding/json"
	"fmt"
	pb "github.com/hazward/plexcluster/plexcluster"
	"google.golang.org/grpc"
	"github.com/streadway/amqp"
	"log"
	"os/exec"
	"time"
)


func handleJobRequest(body []byte, notificationQueue string, channel *amqp.Channel) {
		var job types.Job
		err := json.Unmarshal(body, &job)
		if err != nil {
			log.Println(err)
			return
		}
		if time.Now().After(job.Expiry) {
			log.Println("Job '%s' discard because it expired", job.ID)
			return
		}
		log.Println("Job '%s' received successfully", job.ID)
		runJob(job)

	notification := fmt.Sprintf("Job: %s", job.ID)
	err = channel.Publish(
		"",     // exchange
		notificationQueue, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(notification),
		})
	log.Printf(" [x] Sent nofitification for job '%s'", job.ID)
	if err != nil {
		log.Printf("failed to publish message '%s': %s", notification, err)
	}

}

func runJob(job types.Job) {
	log.Printf("Executing job '%s': %s", job.ID, job.Args)
	cmd := exec.Command("/usr/lib/plexmediaserver/plex_transcoder", job.Args...)
	log.Printf("Running command and waiting for it to finish...")
	err := cmd.Run()
	if err != nil {
		log.Printf("Job '%s' finished with error: %v", job.ID, err)
		return
	}
	log.Printf("Job '%s' finished successfully", job.ID)
}

// Run registers a transcoder and waits to receive jobs from job queue
func Run(transcodeServer string, machineType pb.MachineType) {
	connection, err := grpc.Dial(transcodeServer)
	if err != nil {
		log.Fatalln(err)
	}
	c := pb.NewTranscoderServiceClient(connection)

	c.Transcode()
}
