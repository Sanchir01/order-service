package order

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type Service struct {
	repo      *Repository
	log       *slog.Logger
	primarydb *pgxpool.Pool
}

func NewService(repo *Repository, db *pgxpool.Pool, l *slog.Logger) *Service {
	return &Service{
		repo:      repo,
		log:       l,
		primarydb: db,
	}
}

func (s *Service) CreateOrderService(ctx context.Context,
	props CreateOrderProps,
	paymentprosp CreatePaymentProps,
	deliveryprops CreateDeliveryProps,
	itemsid []uuid.UUID) (interface{}, error) {
	const op = "Wallet.Service.GetBalance"
	log := s.log.With(slog.String("op", op))

	conn, err := s.primarydb.Acquire(ctx)
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
				err := errors.Join(err, rollbackErr)
				log.Error("rollback error", err.Error())
				return
			}
		}

	}()
	orderid, err := s.repo.CreateOrder(ctx, props, tx)
	if err != nil {
		log.Error("CreateOrder error", err.Error())
		return nil, err
	}
	if err := s.repo.CreatePayment(ctx, *orderid, paymentprosp, tx); err != nil {
		log.Error("CreatePayment error", err.Error())
		return nil, err
	}
	if err := s.repo.CreateDelivery(ctx, *orderid, deliveryprops, tx); err != nil {
		log.Error("CreateDelivery error", err.Error())
		return nil, err
	}

	_, err = s.repo.CreateOrderItems(ctx, *orderid, itemsid, tx)
	if err != nil {
		s.log.Error("Order items error", err.Error())
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		s.log.Error("Commit error", err.Error())
		return nil, err
	}
	return nil, nil
}
