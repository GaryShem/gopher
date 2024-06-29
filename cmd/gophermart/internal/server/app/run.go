package app

import (
	"fmt"
	"net/http"

	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/accrual"
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/config"
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/logging"
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/router"
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/memory"
)

func RunServer(sc config.ServerConfig) error {
	if err := logging.InitializeZapLogger("Info"); err != nil {
		return err
	}
	accrualTracker := accrual.NewBonusTracker(sc.AccrualAddress)
	repo, err := memory.NewRepoMemory("", *accrualTracker)
	//repo, err := postgresql.NewRepoPostgreSQL(sc.DBString)
	if err != nil {
		return err
	}
	gRouter, err := router.GopherRouter(
		repo,
		//localMiddleware.LogBody,
	)
	if err != nil {
		return fmt.Errorf("failed to init router: %w", err)
	}
	logging.Log.Infoln("Starting server on address", sc.RunAddress)
	if err = http.ListenAndServe(sc.RunAddress, gRouter); err != nil {
		return err
	}

	return nil
}
