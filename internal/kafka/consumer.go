package kafka

import (
	"context"
	"encoding/json"
	"log/slog"

	"notification-service/entity"
	"notification-service/internal/service"
	"github.com/segmentio/kafka-go"
)

type NotificationConsumer struct {
	reader  *kafka.Reader
	service *service.NotificationService
}

func NewNotificationConsumer(brokers []string, groupID string, service *service.NotificationService) *NotificationConsumer {
	return &NotificationConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   "notifications",
			GroupID: groupID,
		}),
		service: service,
	}
}

// Consume запускается в отдельной горутине
func (c *NotificationConsumer) Consume(ctx context.Context) error {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return err // контекст отменён
		}

		var n model.Notification
		if err := json.Unmarshal(msg.Value, &n); err != nil {
			slog.Error("kafka unmarshal", "error", err)
			continue
		}

		slog.Info("received notification", "id", n.ID)

		if _, err := c.service.SaveNotification(ctx, &n); err != nil {
			slog.Error("save notification", "error", err)
			// можно отправить в DLQ или зафиксировать offset вручную
			continue
		}
	}
}

func (c *NotificationConsumer) Close() error {
	return c.reader.Close()
}