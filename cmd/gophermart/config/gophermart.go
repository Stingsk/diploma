package config

import (
	"errors"

	"github.com/caarlos0/env/v6"
	"github.com/spf13/cobra"
)

const (
	defaultServerAddress      = "127.0.0.1:8080"
	defaultAccrualAddress     = "127.0.0.1:8081"
	defaultDataBaseConnection = "postgresql://localhost:5432/gophermart?sslmode=disable"
)

var (
	ErrInvalidParam = errors.New("invalid param specified")
	rootCmd         = &cobra.Command{
		Use:   "Gophermart",
		Short: "Config for Gophermart",
		Long:  "Config for Gophermart",
	}
	ServerAddress  string
	DatabaseURI    string
	AccrualAddress string
	LogLevel       string
)

type Config struct {
	ServerAddress  string `env:"RUN_ADDRESS"`
	DatabaseURI    string `env:"DATABASE_URI"`
	LogLevel       string `env:"LOG_LEVEL"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func setConfig() {
	rootCmd.Flags().StringVarP(&ServerAddress, "address", "a", defaultServerAddress,
		"Pair of ip:port to listen on")

	rootCmd.Flags().StringVarP(&DatabaseURI, "databaseURI", "d", defaultDataBaseConnection,
		"Database URI for loyalty store")

	rootCmd.Flags().StringVarP(&AccrualAddress, "accrualAddress", "r", defaultAccrualAddress,
		"Pair of ip:port to listen on")

	rootCmd.Flags().StringVarP(&LogLevel, "log-level", "l", "INFO",
		"Set log level: DEBUG|INFO|WARNING|ERROR")
}

func GetGophermartConfig() (Config, error) {
	setConfig()
	LoyaltyServerConfig := Config{
		ServerAddress:  ServerAddress,
		LogLevel:       LogLevel,
		DatabaseURI:    DatabaseURI,
		AccrualAddress: AccrualAddress,
	}
	if err := env.Parse(&LoyaltyServerConfig); err != nil {
		return Config{}, err
	}

	return LoyaltyServerConfig, nil
}
