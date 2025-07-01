package httphandlers

import (
	"log/slog"
	"net/http"

	"github.com/Sanchir01/order-service/internal/app"
	"github.com/Sanchir01/order-service/internal/http/customiddleware"
	"github.com/Sanchir01/order-service/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func StartHTTTPHandlers(handlers *app.Handlers, domain string, l *slog.Logger) http.Handler {
	router := chi.NewRouter()
	custommiddleware(router, l)
	router.Route("/api/v1", func(r chi.Router) {
		r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("helloworl"))
		})
	})
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
	return router
}

func custommiddleware(router *chi.Mux, l *slog.Logger) {
	router.Use(middleware.RequestID, middleware.Recoverer)
	router.Use(middleware.RealIP)
	router.Use(logger.NewMiddlewareLogger(l))
	router.Use(customiddleware.PrometheusMiddleware)
}

func StartPrometheusHandlers() http.Handler {
	router := chi.NewRouter()
	router.Handle("/metrics", promhttp.Handler())
	return router
}
