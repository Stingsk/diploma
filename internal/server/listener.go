package server

import (
	"compress/gzip"
	"net/http"

	"github.com/Stingsk/diploma/internal/logs"
	"github.com/Stingsk/diploma/internal/server/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

func (s *LoyaltyServer) startListener() {
	mux := chi.NewRouter()

	mux.Use(logs.NewStructuredLogger(logrus.StandardLogger()))
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Recoverer)

	compressor := middleware.NewCompressor(gzip.BestCompression)
	mux.Use(compressor.Handler)

	handlers.RegisterPublicHandlers(mux, s.Cfg.UserStore, s.AuthToken())
	handlers.RegisterPrivateHandlers(mux, s.Cfg.OrdersStore, s.AuthToken())

	httpServer := &http.Server{
		Addr:    s.Cfg.ServerAddress,
		Handler: mux,
	}

	s.server = httpServer

	logrus.Info(s.server.ListenAndServe())
}

func (s *LoyaltyServer) stopListener() {
	err := s.server.Shutdown(s.context)
	if err != nil {
		logrus.Info("HTTP server ListenAndServe shut down:", err)
	}
}
