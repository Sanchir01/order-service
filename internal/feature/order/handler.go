package order

import (
	"context"
	"github.com/Sanchir01/order-service/internal/domain/models"
	"github.com/Sanchir01/order-service/pkg/logger"
	"github.com/Sanchir01/order-service/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type OrderHandlers interface {
	GetOrderByIdService(ctx context.Context, orderid uuid.UUID) (*models.OrderFull, error)
	CreateOrderService(ctx context.Context,
		props CreateOrderProps,
		paymentprosp CreatePaymentProps,
		deliveryprops CreateDeliveryProps,
		itemsid []uuid.UUID,
	) error
}
type Handler struct {
	s   OrderHandlers
	log *slog.Logger
}

func NewHandler(s OrderHandlers, l *slog.Logger) *Handler {
	return &Handler{
		s:   s,
		log: l,
	}
}

// @Summary GetOrderById
// @Tags order
// @Description get order by id
// @Accept json
// @Produce json
// @Param id  path string true "order id"
// @Success 200 {object}  GetOrderByIdResponse
// @Failure 400,404 {object}  utils.Response
// @Failure 500 {object}  utils.Response
// @Router /order/id [get]
func (h *Handler) GetOrderById(w http.ResponseWriter, r *http.Request) {
	const op = "Wallet.Handler.GetAllCurrency"
	log := h.log.With(slog.String("op", op))
	id := chi.URLParam(r, "id")
	orderid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed to parse product uuid", logger.Err(err))
		render.JSON(w, r, utils.Error("failed buy product"))
		return
	}
	fullorder, err := h.s.GetOrderByIdService(r.Context(), orderid)
	if err != nil {
		log.Error("failed to get order", logger.Err(err))
		render.JSON(w, r, utils.Error("failed sell order"))
		return
	}
	render.JSON(w, r, GetOrderByIdResponse{
		Response: utils.OK(),
		Data:     *fullorder,
	})
}
