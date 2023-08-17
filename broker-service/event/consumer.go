package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Comsumer struct {
	conn      *amqp.Connection
	queueName string
}

func NewComsumer(conn *amqp.Connection, queueName string) (Comsumer, error) {
	consumer := Comsumer{
		conn:      conn,
		queueName: queueName,
	}
	err := consumer.setup()
	if err != nil {
		return Comsumer{}, err
	}
	return consumer, nil
}

func (c *Comsumer) setup() error {
	channel, err := c.conn.Channel()
	if err != nil {
		return err
	}
	return declareExchange(channel)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (comsumer *Comsumer) Listen(topic []string) error {
	ch, err := comsumer.conn.Channel()
	if err != nil {
		return err
	}
	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}
	for _, s := range topic {
		ch.QueueBind(
			q.Name,
			s,
			"logs_topic",
			false,
			nil,
		)
	}
	if err != nil {
		return err
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range messages {
			var payload Payload
			json.Unmarshal(d.Body, &payload)
			log.Printf("Received a message: %s", payload.Name)
			go handlePayload(payload)
		}
	}()
	fmt.Printf("Waiting for messages [Exchange, Queue] [log_topic, %s]\n", q.Name)
	<-forever
	return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	case "auth":
		// authenticate
		//err := authenticate()
	default:
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}

	}

}

func logEvent(payload Payload) error {
	jsonData, _ := json.MarshalIndent(payload, "", "\t")
	logServiceURL := "http://logger-service/log"
	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusAccepted {
		return err
	}
	return nil
}
