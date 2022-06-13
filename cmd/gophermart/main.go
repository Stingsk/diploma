package main

import (
	"context"
	"log"

	"github.com/Stingsk/diploma/cmd/gophermart/config"
	"github.com/Stingsk/diploma/internal/logs"
	"github.com/Stingsk/diploma/internal/repository"
	"github.com/Stingsk/diploma/internal/server"
	"github.com/sirupsen/logrus"
)

func main() {
	logs.Init()
	cfg, err := config.GetGophermartConfig()
	if err != nil {
		logrus.Fatal("fail to read config server: ", err)
	}
	logrus.Info("config server : ", cfg)
	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logrus.Error("fail to read log level: ", err)
	}
	logrus.SetLevel(level)

	m, err := repository.RunMigration(cfg.DatabaseURI)
	if err != nil && !m {
		log.Fatal(err)
	}

	LoyaltyServerConfig := server.Config{
		ServerAddress:  cfg.ServerAddress,
		DatabaseURI:    cfg.DatabaseURI,
		AccrualAddress: cfg.AccrualAddress,
	}

	loyaltyServer := server.LoyaltyServer{Cfg: &LoyaltyServerConfig}

	loyaltyServer.Run(context.Background())
}
