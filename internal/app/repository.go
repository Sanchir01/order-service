package app

import (
	"github.com/Sanchir01/order-service/internal/feature/events"
	"github.com/Sanchir01/order-service/internal/feature/order"
	db "github.com/Sanchir01/order-service/pkg/database"
	"log/slog"
)

type Repositories struct {
	OrderRepository *order.Repository
	EventRepository *events.Repository
}

func NewRepositories(databases *db.Database, l *slog.Logger) *Repositories {
	return &Repositories{
		OrderRepository: order.NewRepository(databases.PrimaryDB, databases.RedisDB, l),
		EventRepository: events.NewRepository(databases.PrimaryDB),
	}
}
