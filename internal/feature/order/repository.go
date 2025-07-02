package order

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/Sanchir01/order-service/internal/domain/models"
	"github.com/Sanchir01/order-service/pkg/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type Repository struct {
	primarydb *pgxpool.Pool
}

func NewRepository(primarydb *pgxpool.Pool, l *slog.Logger) *Repository {
	return &Repository{primarydb: primarydb}
}

func (r *Repository) GetOrderById(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	query, args, err := sq.Select(
		"o.id AS order_uid",
		"o.track_number",
		"o.entry",
		"o.locale",
		"o.internal_signature",
		"o.customer_id",
		"o.delivery_service",
		"o.shardkey",
		"o.sm_id",
		"o.date_created",
		"o.oof_shard",
		
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
	).
		From("orders o").
		LeftJoin("delivery d ON d.order_uid = o.id").
		LeftJoin("payment p ON p.order_uid = o.id").
		Where(sq.Eq{"o.order_uid": "b563feb7b2b84b6test"}).
		ToSql()
	if err != nil {
		return nil, utils.ErrorQueryString
	}
	return nil, err
}
