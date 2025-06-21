package messagebus

import (
	"context"
	"encoding/json"
	"log"

	"github.com/rabbitmq/amqp091-go"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/pkg/event"
)

type Publisher interface {
	Publish(ctx context.Context, routingKey string, e event.Message) error
}

type Subscriber interface {
	Subscribe(ctx context.Context, routingKey, queueName string, handler func(e event.Message)) error
}

type rabbitmqPublisher struct {
	conn *amqp091.Connection
}

func NewRabbitmqPublisher(conn *amqp091.Connection) Publisher {
	return &rabbitmqPublisher{conn: conn}
}

func (p *rabbitmqPublisher) Publish(ctx context.Context, routingKey string, e event.Message) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	ev, err := json.Marshal(e)
	if err != nil {
		return err
	}

	msg := amqp091.Publishing{
		ContentType: "application/json",
		Body:        ev,
	}

	if err := ch.PublishWithContext(ctx, "amq.topic", routingKey, false, false, msg); err != nil {
		return err
	}

	return nil
}

type rabbitmqSubscriber struct {
	conn *amqp091.Connection
}

func NewRabbitmqSubscriber(conn *amqp091.Connection) Subscriber {
	return &rabbitmqSubscriber{conn: conn}
}

func (p *rabbitmqSubscriber) Subscribe(ctx context.Context, routingKey, queueName string, handler func(e event.Message)) error {
	ch, err := p.conn.Channel()
	if err != nil {
		log.Printf("Failed to create channel: %v", err)
		return err
	}
	msgs, err := ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		log.Printf("Failed to consume messages: %v", err)
		return err
	}

	go func() {
		for d := range msgs {
			var e event.Message
			if err := json.Unmarshal(d.Body, &e); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}
			handler(e)
		}
	}()

	return nil
}
