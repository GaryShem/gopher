package main

import (
	"log"

	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/app"
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/config"
)

func main() {
	sc, err := config.ParseServerFlags()
	if err != nil {
		log.Fatal(err)
	}
	if err = app.RunServer(sc); err != nil {
		log.Fatal(err)
	}
}
