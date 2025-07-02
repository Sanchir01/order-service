package app

import "github.com/Sanchir01/order-service/internal/feature/order"

type Handlers struct {
	OrderHandler *order.Handler
}

func NewHandlers(s *Services) *Handlers {
	return &Handlers{
		OrderHandler: order.NewHandler(s.OrderService),
	}
}
