package handlers

import "github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/repository"

type LoyaltyHandler struct {
	repo repository.Repository
}

func NewLoyaltyHandler(repo repository.Repository) *LoyaltyHandler {
	return &LoyaltyHandler{repo: repo}
}
