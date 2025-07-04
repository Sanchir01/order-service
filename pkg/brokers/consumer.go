package kafkaclient

import (
	"context"
	"encoding/json"
	"github.com/Sanchir01/order-service/internal/feature/order"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"log"
	"log/slog"
)

type OrderService interface {
	CreateOrderService(ctx context.Context,
		props order.CreateOrderProps,
		paymentprosp order.CreatePaymentProps,
		deliveryprops order.CreateDeliveryProps,
		itemsid []uuid.UUID,
	) error
}
type KafkaConsumer struct {
	reader       *kafka.Reader
	log          *slog.Logger
	orderservice OrderService
}

func NewConsumer(topic, broker, groupId string, log *slog.Logger, orderservice OrderService) (*KafkaConsumer, error) {
	cfg := kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: groupId,
	}
	reader := kafka.NewReader(cfg)
	return &KafkaConsumer{
		reader:       reader,
		log:          log,
		orderservice: orderservice,
	}, nil
}
func (kc *KafkaConsumer) Consume(ctx context.Context) error {
	defer kc.reader.Close()
	for {
		select {
		case <-ctx.Done():
			kc.log.Error("🛑 Kafka consumer shutdown")
			return nil
		default:
			msg, err := kc.reader.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return nil
				}
				kc.log.Error("❌ fetch error: %v", err)
				continue
			}
			var ordermsf order.FullOrderMessage
			if err := json.Unmarshal(msg.Value, &ordermsf); err != nil {
				kc.log.Error("❌ failed to unmarshal kafka message", "error", err)
				continue
			}
			kc.log.Info("data consumer", string(msg.Value))
			if err := kc.orderservice.CreateOrderService(ctx, ordermsf.Order, ordermsf.Payment, ordermsf.Delivery, ordermsf.ItemsIds); err != nil {
				kc.log.Error("❌ failed to unmarshal kafka message", "error", err)
				continue
			}
			if err := kc.reader.CommitMessages(ctx, msg); err != nil {
				log.Printf("⚠️ commit error: %v", err)
			}
		}
	}
}
