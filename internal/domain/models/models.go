package models

import (
	"github.com/google/uuid"
	"time"
)

type Delivery struct {
	ID       uuid.UUID `json:"-"`
	OrderUID uuid.UUID `json:"order_uid"`
	Name     string    `json:"name"`
	Phone    string    `json:"phone"`
	Zip      int       `json:"zip"`
	City     string    `json:"city"`
	Address  string    `json:"address"`
	Region   string    `json:"region"`
	Email    string    `json:"email"`
}
type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDT    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}
type Item struct {
	ID          uuid.UUID `json:"id" db:"id"`
	TrackNumber string    `json:"track_number" db:"track_number"`
	Price       int       `json:"price" db:"price"`
	Name        string    `json:"name" db:"name"`
	Sale        string    `json:"sale" db:"sale"`
	Size        int       `json:"size" db:"size"`
	TotalPrice  int       `json:"total_price" db:"total_price"`
	NmID        int       `json:"nm_id" db:"nm_id"`
	Brand       string    `json:"brand" db:"brand"`
	Status      int       `json:"status" db:"status"`
}
type Order struct {
	ID                string    `json:"id"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature,omitempty"`
	CustomerID        string    `json:"customer_id,omitempty"`
	DeliveryService   string    `json:"delivery_service,omitempty"`
	ShardKey          int       `json:"shardkey,omitempty"`
	SmID              int       `json:"sm_id,omitempty"`
	DateCreated       time.Time `json:"date_created"`
}

type OrderFull struct {
	OrderUID          string    `json:"id"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	ShardKey          string    `json:"shardkey"`
	SmID              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
}
