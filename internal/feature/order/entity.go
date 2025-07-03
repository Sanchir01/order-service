package order

import (
	"github.com/Sanchir01/order-service/internal/domain/models"
	"github.com/google/uuid"
	"time"
)

type CreateOrderProps struct {
	TrackNumber       string `json:"track_number"`
	Entry             string `json:"entry"`
	Locale            string `json:"locale"`
	InternalSignature string `json:"internal_signature,omitempty"`
	CustomerID        string `json:"customer_id,omitempty"`
	DeliveryService   string `json:"delivery_service,omitempty"`
	ShardKey          int    `json:"shardkey,omitempty"`
	SmID              int    `json:"sm_id,omitempty"`
}

type CreatePaymentProps struct {
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int64  `json:"amount"`
	PaymentDT    int64  `json:"payment_dt"`
	Region       string `json:"region"`
	Bank         string `json:"bank"`
	DeliveryCost int64  `json:"delivery_cost"`
	GoodsTotal   int64  `json:"goods_total"`
	CustomFree   int64  `json:"custom_free"`
}

type CreateDeliveryProps struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     int    `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}
type ConnectOrderItemsProps struct {
	Items []models.Item `json:"items"`
}
type OrderDB struct {
	ID                string    `db:"id"`
	TrackNumber       string    `db:"track_number"`
	Entry             string    `db:"entry"`
	Locale            string    `db:"locale"`
	InternalSignature string    `db:"internal_signature"`
	CustomerID        string    `db:"customer_id"`
	DeliveryService   string    `db:"delivery_service"`
	ShardKey          int       `db:"shardkey"`
	SmID              int       `db:"sm_id"`
	DateCreated       time.Time `db:"date_created"`
}
type ItemDB struct {
	ID          uuid.UUID `db:"id"`
	TrackNumber string    `db:"track_number"`
	Price       int       `db:"price"`
	Name        string    `db:"name"`
	Sale        string    `db:"sale"`
	Size        int       `db:"size"`
	TotalPrice  int       `db:"total_price"`
	NmID        int       `db:"nm_id"`
	Brand       string    `db:"brand"`
	Status      int       `db:"status"`
}
