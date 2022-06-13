package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Stingsk/diploma/internal/repository/orders"
	"github.com/Stingsk/diploma/internal/repository/users"
	"github.com/go-chi/jwtauth/v5"
	"github.com/sirupsen/logrus"
)

type Config struct {
	ServerAddress  string
	DatabaseURI    string
	AccrualAddress string

	SecretKey []byte

	UserStore   users.Store
	OrdersStore orders.Store

	jwtToken *jwtauth.JWTAuth
}

type LoyaltyServer struct {
	Cfg     *Config
	context context.Context
	server  *http.Server
}

func (s *LoyaltyServer) Run(ctx context.Context) {
	serverContext, serverCancel := context.WithCancel(ctx)
	s.context = serverContext

	s.Cfg.SecretKey = getRandomSecretKey()

	s.Cfg.jwtToken = jwtauth.New("HS256", s.Cfg.SecretKey, s.Cfg.SecretKey)

	closeUsersStore, closeOrdersStore := initStore(s.Cfg)

	pollWorker := PollerWorker{Cfg: PollerConfig{
		AccrualAddress: s.Cfg.AccrualAddress,
		PollInterval:   1000,
	}}

	pollContext, cancelPoller := context.WithCancel(ctx)
	go pollWorker.Run(pollContext, s.Cfg.OrdersStore)

	go s.startListener()
	logrus.Info("start server on " + s.Cfg.ServerAddress)
	<-getSignalChannel()
	logrus.Info("signal received, graceful shutdown the server")
	cancelPoller()
	s.stopListener()

	if err := closeUsersStore(); err != nil {
		logrus.Info("some error occurred while users store close")
	}
	if err := closeOrdersStore(); err != nil {
		logrus.Info("some error occurred while orders store close")
	}

	serverCancel()
}

func (s *LoyaltyServer) AuthToken() *jwtauth.JWTAuth {
	return s.Cfg.jwtToken
}

func getSignalChannel() chan os.Signal {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	return signalChannel
}
