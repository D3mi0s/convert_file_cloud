package queue

import (
	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Conn *amqp091.Connection
	Ch   *amqp091.Channel
}

func NewRabbitMQ(uri string) (*RabbitMQ, error) {
	conn, err := amqp091.Dial(uri)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{Conn: conn, Ch: ch}, nil
}

func (q *RabbitMQ) Publish(queueName, message string) error {
	_, err := q.Ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return q.Ch.Publish(
		"",
		queueName,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}
