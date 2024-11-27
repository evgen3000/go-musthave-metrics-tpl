package router

import (
	"evgen3000/go-musthave-metrics-tpl.git/cmd/server/handlers"
	"evgen3000/go-musthave-metrics-tpl.git/cmd/server/limitedmiddleware"
	"evgen3000/go-musthave-metrics-tpl.git/cmd/server/storage"
	"evgen3000/go-musthave-metrics-tpl.git/internal/compressor"
	"evgen3000/go-musthave-metrics-tpl.git/internal/crypto"
	"evgen3000/go-musthave-metrics-tpl.git/internal/logger"
	"evgen3000/go-musthave-metrics-tpl.git/internal/workerpool"
	"github.com/go-chi/chi/v5"
)

func SetupRouter(storage storage.MetricsStorage, key string, workerpool *workerpool.WorkerPool) *chi.Mux {
	h := handlers.NewHandler(storage, workerpool)
	c := crypto.Crypto{Key: key}
	chiRouter := chi.NewRouter()
	chiRouter.Use(limitedmiddleware.LimitedHandler)
	chiRouter.Use(compressor.GzipMiddleware)
	chiRouter.Use(logger.LoggingMiddleware)

	chiRouter.With(c.HashValidationMiddleware).Route("/update", func(r chi.Router) {
		r.Post("/", h.UpdateMetricHandlerJSON)
		r.Post("/{metricType}/{metricName}/{metricValue}", h.UpdateMetricHandlerText)
	})

	chiRouter.With(c.HashValidationMiddleware).Post("/updates/", h.UpdateMetrics)

	chiRouter.Route("/value", func(r chi.Router) {
		r.Post("/", h.GetMetricHandlerJSON)
		r.Get("/{metricType}/{metricName}", h.GetMetricHandlerText)
	})

	chiRouter.Get("/", h.HomeHandler)

	chiRouter.Get("/ping", h.Ping)

	return chiRouter
}
