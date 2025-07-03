package app

import (
	"context"
	"github.com/Sanchir01/order-service/internal/feature/order"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"log/slog"

	"github.com/Sanchir01/order-service/internal/config"
	db "github.com/Sanchir01/order-service/pkg/database"
	"github.com/Sanchir01/order-service/pkg/logger"
	httpserver "github.com/Sanchir01/order-service/pkg/server/http"
)

type App struct {
	Cfg           *config.Config
	Lg            *slog.Logger
	HttpSrv       *httpserver.Server
	PrometheusSrv *httpserver.Server
	DB            *db.Database
	Handlers      *Handlers
}

func NewApp(ctx context.Context) (*App, error) {
	cfg := config.IntinConfig()
	l := logger.SetupLogger(cfg.Env)
	l.Info("cfg", cfg.Env)
	database, err := db.NewDataBases(cfg, ctx)
	if err != nil {
		l.Error("database", "error", err)
		return nil, err
	}
	httpsrv := httpserver.NewHTTPServer(cfg.HTTPServer.Host, cfg.HTTPServer.Port,
		cfg.HTTPServer.Timeout, cfg.HTTPServer.IdleTimeout)
	prometheusserver := httpserver.NewHTTPServer(cfg.Prometheus.Host, cfg.Prometheus.Port, cfg.Prometheus.Timeout,
		cfg.Prometheus.IdleTimeout)

	repo := NewRepositories(database, l)
	service := NewServices(repo, database, l)
	handler := NewHandlers(service)
	orderprops := order.CreateOrderProps{}
	paymentprops := order.CreatePaymentProps{}
	deliveryprops := order.CreateDeliveryProps{}
	if err := faker.FakeData(&orderprops); err != nil {
		return nil, err
	}
	if err := faker.FakeData(&paymentprops); err != nil {
		return nil, err
	}
	if err := faker.FakeData(&deliveryprops); err != nil {
		return nil, err
	}
	ids := []string{"dfcb6e49-1142-4eec-95ae-d67128deedd3"}
	uuids := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		u, err := uuid.Parse(id)
		if err != nil {
			l.Error("uuid failde parse", "error", err.Error())
			return nil, err
		}
		uuids = append(uuids, u)
	}
	if _, err := service.OrderService.CreateOrderService(ctx, orderprops, paymentprops, deliveryprops, uuids); err != nil {
		l.Error("service order", "error", err.Error())
	}
	idStr := "16a6ad63-1f3c-450c-86a5-74e0b2893da4"
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}
	orderdata, err := repo.OrderRepository.GetOrderById(ctx, id)
	if err != nil {
		l.Error("service order", "error", err.Error())
		return nil, err
	}
	l.Info("order data app", orderdata)
	app := &App{
		Cfg:           cfg,
		Lg:            l,
		HttpSrv:       httpsrv,
		PrometheusSrv: prometheusserver,
		DB:            database,
		Handlers:      handler,
	}
	return app, nil
}
