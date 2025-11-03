package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/kuzin57/auth-system/internal/config"
	"github.com/kuzin57/auth-system/internal/models"
	"github.com/kuzin57/auth-system/internal/services/email"
	"github.com/streadway/amqp"
	"go.uber.org/fx"
)

type MessageBroker struct {
	conn         *amqp.Connection
	ch           *amqp.Channel
	emailService *email.Service
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewMessageBroker(lc fx.Lifecycle, config *config.Config, emailService *email.Service) (*MessageBroker, error) {
	connString := fmt.Sprintf(
		"amqp://%s:%s@%s:%d/%s",
		config.RabbitMQ.User,
		config.RabbitMQ.Password,
		config.RabbitMQ.Host,
		config.RabbitMQ.Port,
		config.RabbitMQ.VHost,
	)

	conn, err := amqp.Dial(connString)
	if err != nil {
		return nil, err
	}

	broker := &MessageBroker{conn: conn, emailService: emailService}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return broker.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return broker.Stop(ctx)
		},
	})

	return broker, nil
}

func (b *MessageBroker) Start(ctx context.Context) error {
	ch, err := b.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	b.ch = ch

	_, err = b.ch.QueueDeclare(
		sendRegistrationLinkQueue,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Создаем отдельный контекст для горутины, который будет отменяться только при остановке
	b.ctx, b.cancel = context.WithCancel(context.Background())

	err = b.StartConsumingSendRegistrationLinkMessage(b.ctx, func(msg models.SendRegistrationLinkMessage) error {
		return b.emailService.HandleSendRegistrationLink(b.ctx, msg)
	})
	if err != nil {
		return fmt.Errorf("failed to start consuming message: %w", err)
	}

	return nil
}

func (b *MessageBroker) Stop(ctx context.Context) error {
	// Отменяем контекст горутины, чтобы она завершилась
	if b.cancel != nil {
		b.cancel()
	}

	if b.ch != nil {
		err := b.ch.Close()
		if err != nil {
			return fmt.Errorf("failed to close channel: %w", err)
		}
	}

	err := b.conn.Close()
	if err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}

	return nil
}

func (b *MessageBroker) PublishSendRegistrationLinkMessage(ctx context.Context, msg models.SendRegistrationLinkMessage) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return b.ch.Publish(
		"",
		sendRegistrationLinkQueue,
		false,
		false,
		amqp.Publishing{Body: msgBytes},
	)
}

func (b *MessageBroker) StartConsumingSendRegistrationLinkMessage(ctx context.Context, handler func(msg models.SendRegistrationLinkMessage) error) error {
	msgs, err := b.ch.Consume(
		sendRegistrationLinkQueue,
		"",
		true,  // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to consume message: %w", err)
	}

	go func() {
		log.Println("starting to consume send registration link message")

		for {
			select {
			case <-ctx.Done():
				log.Println("context done")

				return
			case delivery, ok := <-msgs:
				if !ok {
					log.Println("channel closed")

					return
				}

				log.Println("received message", string(delivery.Body))

				var msg models.SendRegistrationLinkMessage
				if err := json.Unmarshal(delivery.Body, &msg); err != nil {
					log.Println("failed to unmarshal message", err)

					continue
				}

				if err := handler(msg); err != nil {
					log.Println("failed to handle send registration link message", err)

					continue
				}
			}
		}
	}()

	return nil
}
