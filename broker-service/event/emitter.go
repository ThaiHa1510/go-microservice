package event

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	connection *amqp.Connection
}

func (e *Emitter) setup() error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	return declareExchange(channel)
}

func (e *Emitter) Listen(topic []string) error {
	ch, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	for _, s := range topic {
		ch.QueueBind(
			"logs_topic",
			s,
			"logs_topic",
			false,
			nil,
		)
	}
	return nil
}

func (e *Emitter) Push(eventString string, severty string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	log.Println("Pushing to RabbitMQ")
	err = channel.Publish(
		"logs_topic",
		severty,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(eventString),
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		connection: conn,
	}
	err := emitter.setup()
	if err != nil {
		return Emitter{}, err
	}
	return emitter, nil
}
