package clients

import (
	"encoding/json"
	"fmt"
	"log"
	"proyecto_arqui_soft_2/cursos-api/domain"

	"github.com/streadway/amqp"
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

// Crear una nueva instancia del cliente RabbitMQ
func NewRabbit(config RabbitConfig) Rabbit {
	connection, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", config.Username, config.Password, config.Host, config.Port))
	if err != nil {
		log.Fatalf("Error conectando a RabbitMQ: %v", err)
	}
	channel, err := connection.Channel()
	if err != nil {
		log.Fatalf("Error creando canal de RabbitMQ: %v", err)
	}
	queue, err := channel.QueueDeclare(config.QueueName, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Error declarando la cola en RabbitMQ: %v", err)
	}
	return Rabbit{
		connection: connection,
		channel:    channel,
		queue:      queue,
	}
}

// Publicar mensaje en la cola de RabbitMQ
func (r Rabbit) Publish(curso domain.CursoData) error {
	bytes, err := json.Marshal(curso)
	if err != nil {
		return fmt.Errorf("Error al serializar curso: %w", err)
	}
	err = r.channel.Publish(
		"",
		r.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bytes,
		})
	if err != nil {
		return fmt.Errorf("Error al publicar mensaje en RabbitMQ: %w", err)
	}
	return nil
}

// Cerrar conexiones y canal
func (r Rabbit) Close() {
	if err := r.channel.Close(); err != nil {
		log.Printf("Error cerrando canal de RabbitMQ: %v", err)
	}
	if err := r.connection.Close(); err != nil {
		log.Printf("Error cerrando conexión de RabbitMQ: %v", err)
	}
}
