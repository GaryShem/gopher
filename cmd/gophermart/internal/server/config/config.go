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

func parseCommandLineArguments() ServerConfig {
	result := ServerConfig{}
	flag.StringVar(&result.RunAddress, "a", "localhost:8080", "service address")
	flag.StringVar(&result.DBString, "d", "host=localhost user=postgres password=1231 dbname=postgres sslmode=disable", "database URI")
	flag.StringVar(&result.AccrualAddress, "r", "http://localhost:8081", "accrual address")
	flag.Parse()
	return result
}

func parseEnvArguments(sc ServerConfig) (ServerConfig, error) {
	var ec ServerConfig
	if err := env.Parse(&ec); err != nil {
		return sc, err
	}
	if ec.RunAddress != "" {
		sc.RunAddress = ec.RunAddress
	}
	if ec.DBString != "" {
		sc.DBString = ec.DBString
	}
	if ec.AccrualAddress != "" {
		sc.AccrualAddress = ec.AccrualAddress
	}
	return sc, nil
}

func ParseServerFlags() (ServerConfig, error) {
	configCLI := parseCommandLineArguments()
	return parseEnvArguments(configCLI)
}
