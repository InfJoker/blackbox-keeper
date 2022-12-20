package rabbit

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitWriter struct {
	ch *amqp.Channel
	q  amqp.Queue
}

func NewWriter(url, queue string) (*RabbitWriter, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	q, err := ch.QueueDeclare(
		queue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &RabbitWriter{
		ch: ch,
		q:  q,
	}, nil
}

func (r *RabbitWriter) Write(p []byte) (int, error) {
	err := r.ch.Publish(
		"",
		r.q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "textplain",
			Body:        p,
		},
	)
	if err != nil {
		return 0, nil
	}
	return len(p), nil
}
