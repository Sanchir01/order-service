package app

import (
	"context"
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
	service := NewServices(repo)
	handler := NewHandlers(service)
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
