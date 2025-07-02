package main

import (
	"context"
	"errors"
	"github.com/Sanchir01/order-service/internal/http"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	apps "github.com/Sanchir01/order-service/internal/app"
)

// @title 🚀 Currency Wallet
// @version         1.0
// @description This is a sample server seller
// @termsOfService  http://swagger.io/terms/

// @host localhost:5000
// @BasePath /api/v1

// @securityDefinitions.apikey AccessTokenCookie
// @in cookie
// @name accessToken

// @securityDefinitions.apikey RefreshTokenCookie
// @in cookie
// @name refreshToken

// @contact.name GitHub
// @contact.url https://github.com/Sanchir01
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	app, err := apps.NewApp(ctx)
	if err != nil {
		panic(err)
	}

	app.Lg.Info("started server", app.Cfg.HTTPServer.Port)
	go func() {
		if err := app.HttpSrv.Run(httphandlers.StartHTTTPHandlers(app.Handlers, app.Cfg.Domain, app.Lg)); err != nil {
			if !errors.Is(err, context.Canceled) {
				app.Lg.Error("Listen prometheus server error", slog.String("error", err.Error()))
				return
			}
		}
	}()
	go func() {
		if err := app.PrometheusSrv.Run(httphandlers.StartPrometheusHandlers()); err != nil {
			if !errors.Is(err, context.Canceled) {
				app.Lg.Error("Listen prometheus server error", slog.String("error", err.Error()))
				return
			}
		}
	}()
	<-ctx.Done()
	if err := app.HttpSrv.Gracefull(ctx); err != nil {
		app.Lg.Error("server gracefull")
	}
	if err := app.DB.Close(); err != nil {
		app.Lg.Error("Close database", slog.String("error", err.Error()))
	}
}
