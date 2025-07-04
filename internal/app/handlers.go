package app

import (
	"github.com/Sanchir01/order-service/internal/feature/order"
	"log/slog"
)

type Handlers struct {
	OrderHandler *order.Handler
}

func NewHandlers(s *Services, l *slog.Logger) *Handlers {
	return &Handlers{
		OrderHandler: order.NewHandler(s.OrderService, l),
	}
}
