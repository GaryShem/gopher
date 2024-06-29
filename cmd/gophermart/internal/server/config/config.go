package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type ServerConfig struct {
	RunAddress     string `env:"RUN_ADDRESS"`
	DBString       string `env:"DATABASE_URI"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func ParseServerFlags() (ServerConfig, error) {
	result := ServerConfig{}
	flag.StringVar(&result.RunAddress, "a", "localhost:8080", "service address")
	flag.StringVar(&result.DBString, "d", "", "database URI")
	flag.StringVar(&result.AccrualAddress, "r", "", "accrual address")
	flag.Parse()

	var ec ServerConfig
	if err := env.Parse(&ec); err != nil {
		return result, err
	}
	if ec.RunAddress != "" {
		result.RunAddress = ec.RunAddress
	}
	if ec.DBString != "" {
		result.DBString = ec.DBString
	}
	if ec.AccrualAddress != "" {
		result.AccrualAddress = ec.AccrualAddress
	}
	return result, nil
}
