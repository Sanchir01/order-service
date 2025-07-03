package app

import (
	"github.com/Sanchir01/order-service/internal/feature/order"
	db "github.com/Sanchir01/order-service/pkg/database"
	"log/slog"
)

type Services struct {
	OrderService *order.Service
}

func NewServices(r *Repositories, databases *db.Database, l *slog.Logger) *Services {
	return &Services{
		OrderService: order.NewService(r.OrderRepository, databases.PrimaryDB, l),
	}
}
