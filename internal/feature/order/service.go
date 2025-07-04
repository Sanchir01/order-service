package order

import (
	"context"
	"errors"
	"github.com/Sanchir01/order-service/internal/domain/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type OrderService interface {
	CreateDelivery(ctx context.Context, ordserid uuid.UUID, payment CreateDeliveryProps, tx pgx.Tx) error
	CreatePayment(ctx context.Context, ordserid uuid.UUID, payment CreatePaymentProps, tx pgx.Tx) error
	CreateOrder(ctx context.Context, order CreateOrderProps, tx pgx.Tx) (*uuid.UUID, error)
	CreateOrderItems(
		ctx context.Context, orderID uuid.UUID, productIDs []uuid.UUID, tx pgx.Tx,
	) ([]uuid.UUID, error)
	GetOrderById(ctx context.Context, id uuid.UUID) (*models.OrderFull, error)
}
type Service struct {
	repo      OrderService
	log       *slog.Logger
	primarydb *pgxpool.Pool
}

func NewService(repo OrderService, db *pgxpool.Pool, l *slog.Logger) *Service {
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
	itemsid []uuid.UUID,
) error {
	const op = "Order.Service.CreateOrderService"
	log := s.log.With(slog.String("op", op))

	conn, err := s.primarydb.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Error("tx error", err.Error())
		return err
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
		return err
	}
	if err := s.repo.CreatePayment(ctx, *orderid, paymentprosp, tx); err != nil {
		log.Error("CreatePayment error", err.Error())
		return err
	}
	if err := s.repo.CreateDelivery(ctx, *orderid, deliveryprops, tx); err != nil {
		log.Error("CreateDelivery error", err.Error())
		return err
	}

	_, err = s.repo.CreateOrderItems(ctx, *orderid, itemsid, tx)
	if err != nil {
		log.Error("Order items error", err.Error())
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		s.log.Error("Commit error", err.Error())
		return err
	}
	return nil
}
func (s *Service) GetOrderByIdService(ctx context.Context, orderid uuid.UUID) (*models.OrderFull, error) {
	const op = "Order.Service.CreateOrderService"
	log := s.log.With(slog.String("op", op))
	fullorder, err := s.repo.GetOrderById(ctx, orderid)
	if err != nil {
		log.Error("GetOrderById error", err.Error())
		return nil, err
	}

	return fullorder, nil
}
