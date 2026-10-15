package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const UserDeletedQueue = "user_deleted"

type userDeletedMessage struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

type RabbitMQPublisher struct {
	mu sync.Mutex
	ch *amqp.Channel
}

func NewRabbitMQPublisher(conn *amqp.Connection) (*RabbitMQPublisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	if err := ch.Confirm(false); err != nil {
		return nil, err
	}
	if _, err := ch.QueueDeclare(UserDeletedQueue, true, false, false, false, nil); err != nil {
		return nil, err
	}
	return &RabbitMQPublisher{ch: ch}, nil
}

func (p *RabbitMQPublisher) PublishUserDeleted(ctx context.Context, userID, email string) error {
	payload, err := json.Marshal(userDeletedMessage{UserID: userID, Email: email})
	if err != nil {
		return err
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	confirm, err := p.ch.PublishWithDeferredConfirmWithContext(ctx, "", UserDeletedQueue, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		MessageId:    userID,
		Timestamp:    time.Now(),
		Body:         payload,
	})
	if err != nil {
		return err
	}
	acked, err := confirm.WaitContext(ctx)
	if err != nil {
		return err
	}
	if !acked {
		return fmt.Errorf("broker nacked user_deleted for %s", userID)
	}
	return nil
}
