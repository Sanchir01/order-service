package httphandlers

import (
	_ "github.com/Sanchir01/order-service/docs"
	"github.com/Sanchir01/order-service/internal/app"
	"github.com/Sanchir01/order-service/internal/http/customiddleware"
	"github.com/Sanchir01/order-service/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"log/slog"
	"net/http"
)

func StartHTTTPHandlers(handlers *app.Handlers, domain string, l *slog.Logger) http.Handler {
	router := chi.NewRouter()
	newChiCors(router)
	custommiddleware(router, l)
	router.Route("/api/v1", func(r chi.Router) {
		r.Get("/order/{id}", handlers.OrderHandler.GetOrderById)
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
func newChiCors(r chi.Router) {
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		MaxAge:           300,
	}))

}
func StartPrometheusHandlers() http.Handler {
	router := chi.NewRouter()
	router.Handle("/metrics", promhttp.Handler())
	return router
}
