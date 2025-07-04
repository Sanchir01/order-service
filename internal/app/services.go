package app

import (
	"github.com/Sanchir01/order-service/internal/feature/events"
	"github.com/Sanchir01/order-service/internal/feature/order"
	kafkaclient "github.com/Sanchir01/order-service/pkg/brokers"
	db "github.com/Sanchir01/order-service/pkg/database"
	"log/slog"
)

type Services struct {
	OrderService *order.Service
	EventService *events.Service
}

func NewServices(r *Repositories, databases *db.Database, l *slog.Logger, producer *kafkaclient.Producer) *Services {
	return &Services{
		OrderService: order.NewService(r.OrderRepository, databases.PrimaryDB, l),
		EventService: events.NewEventService(l, r.EventRepository, producer),
	}
}
