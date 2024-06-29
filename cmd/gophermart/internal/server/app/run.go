package app

import (
	"fmt"
	"net/http"

	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/config"
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/logging"
	localMiddleware "github.com/GaryShem/gopher/cmd/gophermart/internal/server/middleware"
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/router"
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/memory"
)

func RunServer(sc config.ServerConfig) error {
	if err := logging.InitializeZapLogger("Info"); err != nil {
		return err
	}
	repo := memory.NewRepoMemory("")
	gRouter, err := router.GopherRouter(repo, localMiddleware.LogBody)
	if err != nil {
		return fmt.Errorf("failed to init router: %w", err)
	}
	logging.Log.Infoln("Starting server on address", sc.RunAddress)
	if err = http.ListenAndServe(sc.RunAddress, gRouter); err != nil {
		return err
	}

	return nil
}
