package queues

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	cursos"proyecto_arqui_soft_2/search-api/domain"
	"proyecto_arqui_soft_2/search-api/services"
)

type RabbitConfig struct {
	Host      string
	Port      string
	Username  string
	Password  string
	QueueName string
}

type Rabbit struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
}

// NewRabbit creates a new RabbitMQ connection and declares the queue
func NewRabbit(config RabbitConfig) Rabbit {
	connection, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", config.Username, config.Password, config.Host, config.Port))
	if err != nil {
		log.Fatalf("error getting Rabbit connection: %w", err)
	}
	channel, err := connection.Channel()
	if err != nil {
		log.Fatalf("error creating Rabbit channel: %w", err)
	}
	queue, err := channel.QueueDeclare(config.QueueName, false, false, false, false, nil)
	return Rabbit{
		connection: connection,
		channel:    channel,
		queue:      queue,
	}
}

// StartConsumer starts listening for messages on the RabbitMQ queue
func (queue Rabbit) StartConsumer(service services.Service) error {
	messages, err := queue.channel.Consume(
		queue.queue.Name,
		"",
		true, // Auto-acknowledge messages
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error registering consumer: %w", err)
	}

	go func() {
		for msg := range messages {
			var CursoNew cursos.CursoNew
			if err := json.Unmarshal(msg.Body, &CursoNew); err != nil {
				log.Printf("error unmarshaling message: %v", err)
				continue
			}
			fmt.Println(CursoNew)
			// Call the service method
			service.HandleCursoNew(CursoNew)


		}
	}()

	return nil
}

// Close cleans up the RabbitMQ resources
func (queue Rabbit) Close() {
	if err := queue.channel.Close(); err != nil {
		log.Printf("error closing Rabbit channel: %v", err)
	}
	if err := queue.connection.Close(); err != nil {
		log.Printf("error closing Rabbit connection: %v", err)
	}
}
