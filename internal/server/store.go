package server

import (
	"database/sql"

	"github.com/Stingsk/diploma/internal/repository/orders"
	"github.com/Stingsk/diploma/internal/repository/users"
	"github.com/sirupsen/logrus"
)

const (
	psqlDriverName = "pgx"
)

func initStore(config *Config) (func() error, func() error) {
	conn, err := sql.Open(psqlDriverName, config.DatabaseURI)
	if err != nil {
		logrus.Error("couldn't create database connection")
	}

	userStore := users.NewDBStore(conn)
	config.UserStore = userStore
	logrus.Info("using database for user storage")

	ordersStore := orders.NewDBStore(conn)
	config.OrdersStore = ordersStore
	logrus.Info("using database for orders storage")

	return func() error {
			return userStore.Close()
		}, func() error {
			return ordersStore.Close()
		}
}
