package publisher

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

// EventPublisher define o contrato
type EventPublisher interface {
	Publish(action string, payload interface{}) error
}

// RabbitPublisher implementa via RabbitMQ
type RabbitPublisher struct {
	ch        *amqp.Channel
	queueName string
}

func NewRabbitPublisher(ch *amqp.Channel, queueName string) *RabbitPublisher {
	return &RabbitPublisher{ch: ch, queueName: queueName}
}

func (p *RabbitPublisher) Publish(action string, payload interface{}) error {
	evt := struct {
		Action     string      `json:"action"`
		SimpleUser interface{} `json:"user"`
	}{action, payload}

	body, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	return p.ch.Publish(
		"", p.queueName,
		false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
