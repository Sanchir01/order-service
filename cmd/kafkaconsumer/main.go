package main

import (
	"context"
	apps "github.com/Sanchir01/order-service/internal/app"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	app, err := apps.NewApp(ctx)
	if err != nil {
		panic(err)
	}
	app.Lg.Info("Starting application...")
	<-ctx.Done()
	app.DB.Close()

}
