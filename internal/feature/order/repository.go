package order

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/Sanchir01/order-service/internal/domain/models"
	"github.com/Sanchir01/order-service/pkg/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"log/slog"
)

type Repository struct {
	primarydb *pgxpool.Pool
	log       *slog.Logger
}

func NewRepository(primarydb *pgxpool.Pool, l *slog.Logger) *Repository {
	return &Repository{primarydb: primarydb, log: l}
}

func (r *Repository) GetOrderById(ctx context.Context, id uuid.UUID) (*models.OrderFull, error) {
	conn, err := r.primarydb.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	query, args, err := sq.Select(
		"o.id",
		"o.track_number",
		"o.entry",
		"o.locale",
		"o.internal_signature",
		"o.customer_id",
		"o.delivery_service",
		"o.shardkey",
		"o.sm_id",
		"o.date_created",

		"d.name AS delivery_name",
		"d.phone AS delivery_phone",
		"d.zip AS delivery_zip",
		"d.city AS delivery_city",
		"d.address AS delivery_address",
		"d.region AS delivery_region",
		"d.email AS delivery_email",

		"p.transaction AS payment_transaction",
		"p.request_id AS payment_request_id",
		"p.currency AS payment_currency",
		"p.provider AS payment_provider",
		"p.amount AS payment_amount",
		"p.payment_dt AS payment_dt",
		"p.bank AS payment_bank",
		"p.delivery_cost AS payment_delivery_cost",
		"p.goods_total AS payment_goods_total",
		"p.custom_fee AS payment_custom_fee",

		"i.id",
		"i.track_number",
		"i.price",
		"i.name",
		"i.sale",
		"i.size",
		"i.total_price",
		"i.nm_id",
		"i.brand",
		"i.status",
	).
		From("orders o").
		LeftJoin("delivery d ON d.order_uid = o.id").
		LeftJoin("payment p ON p.order_uid = o.id").
		LeftJoin("order_items oi ON oi.order_uid = o.id").
		LeftJoin("items i ON i.id = oi.item_id").
		Where(sq.Eq{"o.id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, utils.ErrorQueryString
	}
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var order models.OrderFull
	var items []models.Item
	for rows.Next() {
		var item models.Item
		err := rows.Scan(
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.ShardKey,
			&order.SmID,
			&order.DateCreated,

			&order.Delivery.Name,
			&order.Delivery.Phone,
			&order.Delivery.Zip,
			&order.Delivery.City,
			&order.Delivery.Address,
			&order.Delivery.Region,
			&order.Delivery.Email,

			&order.Payment.Transaction,
			&order.Payment.RequestID,
			&order.Payment.Currency,
			&order.Payment.Provider,
			&order.Payment.Amount,
			&order.Payment.PaymentDT,
			&order.Payment.Bank,
			&order.Payment.DeliveryCost,
			&order.Payment.GoodsTotal,
			&order.Payment.CustomFee,

			&item.ID,
			&item.TrackNumber,
			&item.Price,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			r.log.Error("fail to scan  order", "err", err.Error())
			return nil, err
		}
		items = append(items, item)
	}

	order.Items = items
	return &order, nil
}

func (r *Repository) CreateOrder(ctx context.Context, order CreateOrderProps, tx pgx.Tx) (*uuid.UUID, error) {
	query, args, err := sq.
		Insert("orders").
		Columns("track_number", "entry", "locale", "internal_signature", "customer_id",
			"delivery_service", "shardkey", "sm_id",
		).Values(
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.ShardKey,
		order.SmID,
	).Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, utils.ErrorQueryString
	}
	var orderid uuid.UUID
	if err := tx.QueryRow(ctx, query, args...).Scan(&orderid); err != nil {
		return nil, err
	}
	return &orderid, err
}

func (r *Repository) CreatePayment(ctx context.Context, ordserid uuid.UUID, payment CreatePaymentProps, tx pgx.Tx) error {
	query, args, err := sq.
		Insert("payment").
		Columns(
			"order_uid",
			"currency",
			"provider",
			"amount",
			"payment_dt",
			"bank",
			"delivery_cost",
			"goods_total",
			"custom_fee",
		).
		Values(
			ordserid,
			payment.Currency,
			payment.Provider,
			payment.Amount,
			payment.PaymentDT,
			payment.Bank,
			payment.DeliveryCost,
			payment.GoodsTotal,
			payment.CustomFree,
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return utils.ErrorQueryString
	}

	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return err
	}
	return nil
}

func (r *Repository) CreateDelivery(ctx context.Context, ordserid uuid.UUID, payment CreateDeliveryProps, tx pgx.Tx) error {
	query, args, err := sq.Insert("delivery").
		Columns("order_uid", "name", "phone", "zip", "city", "address", "region", "email").
		Values(
			ordserid,
			payment.Name,
			payment.Phone,
			payment.Zip,
			payment.City,
			payment.Address,
			payment.Region,
			payment.Email,
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return utils.ErrorQueryString
	}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return err
	}
	return nil
}

func (r *Repository) CreateOrderItems(
	ctx context.Context, orderID uuid.UUID, productIDs []uuid.UUID, tx pgx.Tx,
) ([]uuid.UUID, error) {
	if len(productIDs) == 0 {
		return []uuid.UUID{}, fmt.Errorf("productIDs is empty")
	}
	queryBuilder := sq.Insert("order_items").Columns("order_uid", "item_id")

	for i := 0; i < len(productIDs); i++ {
		queryBuilder = queryBuilder.Values(orderID, productIDs[i])
	}

	query, args, err := queryBuilder.Suffix("RETURNING id").PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var ids []uuid.UUID
	log.Printf("SQL Query: %s", query)
	log.Printf("Args: %v", args)

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		log.Printf("Error executing query: %v", err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	log.Printf("ids items: %v", ids)
	return ids, nil
}
