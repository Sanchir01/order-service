package main

import (
	"context"
	"encoding/json"
	sq "github.com/Masterminds/squirrel"
	apps "github.com/Sanchir01/order-service/internal/app"
	"github.com/Sanchir01/order-service/internal/feature/order"
	"github.com/Sanchir01/order-service/pkg/utils"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	app, err := apps.NewApp(ctx)
	if err != nil {
		panic(err)
	}
	app.Lg.Info("Starting application...")

	ticker := time.NewTicker(3 * time.Second)
	brokerAddress := "localhost:9092"
	topic := "notification"

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{brokerAddress},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	log.Println("Message sent successfully!")

	go func() {
		defer ticker.Stop()
		defer func() {
			if err := recover(); err != nil {
				app.Lg.Warn("Recovered from panic:", err, "\nStack trace:\n", string(debug.Stack()))
			}
		}()
		for {
			select {
			case <-ctx.Done():
				app.Lg.Error("stopping event service")
				return
			case <-ticker.C:
				ids, err := getAllItemsIds(ctx)
				if err != nil {
					app.Lg.Error("error getting all items ids: ", err)
					continue
				}
				var createorder order.CreateOrderProps
				if err := faker.FakeData(&createorder); err != nil {
					app.Lg.Error(err.Error())
					continue
				}
				var createpayment order.CreatePaymentProps
				if err := faker.FakeData(&createpayment); err != nil {
					app.Lg.Error(err.Error())
					continue
				}
				var createdelivery order.CreateDeliveryProps
				if err := faker.FakeData(&createdelivery); err != nil {
					app.Lg.Error(err.Error())
					continue
				}
				payload := order.FullOrderMessage{
					Order:    createorder,
					Payment:  createpayment,
					Delivery: createdelivery,
					ItemsIds: ids,
				}
				msgBytes, err := json.Marshal(payload)
				if err != nil {
					app.Lg.Error("json marshal error", "err", err)
					continue
				}
				app.Lg.Info("json marshaled", string(msgBytes))
				err = writer.WriteMessages(ctx, kafka.Message{
					Value: msgBytes,
				})
				if err != nil {
					app.Lg.Error("failed to write message", "err", err)
				}
			}
		}
	}()
	<-ctx.Done()

}

func getAllItemsIds(ctx context.Context) ([]uuid.UUID, error) {
	dsn := "postgres://postgres:postgres@localhost:5441/order?sslmode=disable"

	var pool *pgxpool.Pool

	err := utils.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		var err error
		pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			return err
		}
		return nil
	}, 1, 1*time.Second)
	if err != nil {
		return nil, err
	}
	defer pool.Close()

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query, args, err := sq.
		Select("id").
		From("items").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var itemsIDs []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		itemsIDs = append(itemsIDs, id)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return itemsIDs, nil
}
