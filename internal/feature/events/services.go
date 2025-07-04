package events

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Sanchir01/order-service/pkg/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"time"
)

type EventRepositoryInterface interface {
	CreateEvent(ctx context.Context, eventType, payload string, tx pgx.Tx) (uuid.UUID, error)
	GetManyEvents(ctx context.Context, limit uint64) ([]*EventDB, error)
	SetDone(ctx context.Context, ids []uuid.UUID) error
}
type EventSender interface {
	Produce(message string, value []byte) error
}
type Service struct {
	log       *slog.Logger
	primarydb *pgxpool.Pool
	repo      EventRepositoryInterface
	kaf       EventSender
}

func NewEventService(log *slog.Logger, repo EventRepositoryInterface, kaf EventSender, primarydb *pgxpool.Pool) *Service {
	return &Service{
		log,
		primarydb,
		repo,
		kaf,
	}
}

func (e *Service) CreateEvent(ctx context.Context, eventType, payload string) (*uuid.UUID, error) {
	const op = "Wallet.Handler.GetAllCurrency"
	log := e.log.With(slog.String("op", op))

	conn, err := e.primarydb.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Error("tx error", err.Error())
		return nil, err
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				err = errors.Join(err, rollbackErr)
				log.Error("rollback error", err.Error())
			}
		}
	}()

	eventid, err := e.repo.CreateEvent(ctx, eventType, payload, tx)
	if err != nil {
		log.Error("failed create event", err.Error())
		return nil, err
	}
	log.Info("eventid", "id", eventid)

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &eventid, nil
}

func (e *Service) StartCreateEvent(ctx context.Context, handlePeriod time.Duration, limitEvents uint64, topic string) {
	const op = "EventService.StartCreateEvent"

	log := e.log.With(slog.String("op", op))
	ticker := time.NewTicker(handlePeriod)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Info("stopping event service")
				return

			case <-ticker.C:
				log.Debug("starting process events")

				events, err := e.repo.GetManyEvents(ctx, limitEvents)
				if err != nil {
					log.Error("failed to get new events", logger.Err(err))
					continue
				}

				if len(events) == 0 {
					log.Debug("no events to process")
					continue
				}

				var ids []uuid.UUID
				for _, event := range events {
					ids = append(ids, event.ID)
				}
				for _, event := range events {
					if err := e.SendMessage(event, topic); err != nil {
						log.Error("failed to send event", logger.Err(err))
					}
				}

				if err := e.repo.SetDone(ctx, ids); err != nil {
					log.Error("failed to set events done", logger.Err(err))
					continue
				}

				log.Info("successfully processed events", slog.Int("count", len(ids)))
			}
		}
	}()
}

func (e *Service) SendMessage(event *EventDB, topic string) error {
	const op = "services.event-sender.SendMessage"

	log := e.log.With(slog.String("op", op))
	jsondata, err := json.Marshal(event)
	if err != nil {
		return err
	}
	if err := e.kaf.Produce(topic, jsondata); err != nil {
		return err
	}
	log.Info("successfully sent message", slog.Int("count", len(jsondata)))
	return err
}
