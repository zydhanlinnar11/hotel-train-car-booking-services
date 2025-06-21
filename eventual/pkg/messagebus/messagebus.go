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

	payload, err := json.Marshal(e.Payload)
	if err != nil {
		return err
	}

	msg := amqp091.Publishing{
		ContentType: "application/json",
		Body:        payload,
	}

	if err := ch.PublishWithContext(ctx, "", routingKey, false, false, msg); err != nil {
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
		return err
	}
	defer ch.Close()

	msgs, err := ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case d, ok := <-msgs:
				if !ok {
					log.Printf("Message channel closed for queue: %s", queueName)
					return
				}
				var e event.Message
				if err := json.Unmarshal(d.Body, &e); err != nil {
					log.Printf("Failed to unmarshal message: %v", err)
					continue
				}
				handler(e)
			case <-ctx.Done():
				log.Printf("Context cancelled, stopping subscriber for queue: %s", queueName)
				return
			}
		}
	}()

	return nil
}
