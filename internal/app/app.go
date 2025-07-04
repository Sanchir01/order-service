package app

import (
	"context"
	"github.com/Sanchir01/order-service/internal/config"
	"github.com/Sanchir01/order-service/internal/profiling"
	kafkaclient "github.com/Sanchir01/order-service/pkg/brokers"
	db "github.com/Sanchir01/order-service/pkg/database"
	"github.com/Sanchir01/order-service/pkg/logger"
	httpserver "github.com/Sanchir01/order-service/pkg/server/http"
	"log/slog"
)

type App struct {
	Cfg           *config.Config
	Lg            *slog.Logger
	HttpSrv       *httpserver.Server
	PrometheusSrv *httpserver.Server
	DB            *db.Database
	Services      *Services
	Handlers      *Handlers
	Consumer      *kafkaclient.KafkaConsumer
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

	kaf, err := kafkaclient.NewProducer(cfg.Kafka.Notification.Broke, cfg.Kafka.Notification.Topic[0], cfg.Kafka.Notification.Retries, ctx)
	if err != nil {
		return nil, err
	}
	repo := NewRepositories(database, l)
	service := NewServices(repo, database, l, kaf)
	handler := NewHandlers(service, l)

	consumer, err := kafkaclient.NewConsumer(cfg.Kafka.Consumer.Topic[0], cfg.Kafka.Consumer.Broke[0], cfg.Kafka.Consumer.GroupId, l, service.OrderService)
	if err := profiling.InitPyroscope(); err != nil {
		return nil, err
	}
	app := &App{
		Cfg:           cfg,
		Lg:            l,
		HttpSrv:       httpsrv,
		PrometheusSrv: prometheusserver,
		DB:            database,
		Handlers:      handler,
		Consumer:      consumer,
	}
	return app, nil
}
