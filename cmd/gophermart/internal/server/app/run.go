package app

import (
	"fmt"
	"net/http"

	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/accrual"
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/config"
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/logging"
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/router"
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/postgresql"
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/repository"
)

func initLogger() error {
	if err := logging.InitializeZapLogger("Info"); err != nil {
		return err
	}
	return nil
}

func initRepo(sc config.ServerConfig) (repository.Repository, error) {
	accrualTracker := accrual.NewBonusTracker(sc.AccrualAddress)
	//repo, err := memory.NewRepoMemory("", *accrualTracker)
	repo, err := postgresql.NewRepoPostgreSQL(sc.DBString, *accrualTracker)
	return repo, err
}

func startRouter(sc config.ServerConfig, repo repository.Repository) error {
	gRouter, err := router.GopherRouter(repo)
	if err != nil {
		return fmt.Errorf("failed to init router: %w", err)
	}
	logging.Log.Infoln("Starting server on address", sc.RunAddress)
	if err = http.ListenAndServe(sc.RunAddress, gRouter); err != nil {
		return err
	}
	return nil
}

func RunServer(sc config.ServerConfig) error {
	if err := initLogger(); err != nil {
		return err
	}
	repo, err := initRepo(sc)
	if err != nil {
		return err
	}
	err = startRouter(sc, repo)
	return err
}
