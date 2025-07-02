package app

import "github.com/Sanchir01/order-service/internal/feature/order"

type Services struct {
	OrderService *order.Service
}

func NewServices(r *Repositories) *Services {
	return &Services{
		OrderService: order.NewService(r.OrderRepository),
	}
}
