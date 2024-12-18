package router

import (
	"evgen3000/go-musthave-metrics-tpl.git/cmd/server/handlers"
	"evgen3000/go-musthave-metrics-tpl.git/cmd/server/storage"
	"evgen3000/go-musthave-metrics-tpl.git/internal/compressor"
	"evgen3000/go-musthave-metrics-tpl.git/internal/logger"
	"github.com/go-chi/chi/v5"
)

func SetupRouter(storage storage.MetricsStorage) *chi.Mux {
	h := handlers.NewHandler(storage)
	chiRouter := chi.NewRouter()
	chiRouter.Use(logger.LoggingMiddleware)
	chiRouter.Use(compressor.GzipMiddleware)

	chiRouter.Route("/update", func(r chi.Router) {
		r.Post("/", h.UpdateMetricHandlerJSON)
		r.Post("/{metricType}/{metricName}/{metricValue}", h.UpdateMetricHandlerText)
	})

	chiRouter.Post("/updates/", h.UpdateMetrics)

	chiRouter.Route("/value", func(r chi.Router) {
		r.Post("/", h.GetMetricHandlerJSON)
		r.Get("/{metricType}/{metricName}", h.GetMetricHandlerText)
	})

	chiRouter.Get("/", h.HomeHandler)

	chiRouter.Get("/ping", h.Ping)

	return chiRouter
}
