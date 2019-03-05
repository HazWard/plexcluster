package queue

import (
	"fmt"
	"github.com/hazward/plexcluster/types"
	"github.com/streadway/amqp"
)

type RabbitMQQueue struct {
	queueConnection *amqp.Connection
	queueChannel *amqp.Channel
	jobSubmissionRoutingKey string
	notificationRoutingKey string
}


func NewRabbitMQQueue(uri string, jobSubmissionQueueName, notificationQueueName string) (*RabbitMQQueue, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	submissionQueue, err := channel.QueueDeclare(
		jobSubmissionQueueName, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	if err != nil {
		return nil, fmt.Errorf("error while declaring job submission queue: %s", err)
	}

	notificationQueue, err := channel.QueueDeclare(
		notificationQueueName, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	if err != nil {
		return nil, fmt.Errorf("error while declaring notification queue: %s", err)
	}

	return &RabbitMQQueue{
		queueConnection: conn,
		queueChannel: channel,
		jobSubmissionRoutingKey: submissionQueue.Name,
		notificationRoutingKey: notificationQueue.Name,
	}, nil
}

func (r *RabbitMQQueue) Submit(job types.Job) error {
	return fmt.Errorf("unimplemented function")
}

func (r *RabbitMQQueue) WaitForCompletion(jobID string, found chan int) error {
	return fmt.Errorf("unimplemented function")
}