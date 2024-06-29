package router

import (
	chi "github.com/go-chi/chi/v5"

	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/handlers"
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/middleware"
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/repository"
)

func GopherRouter(repo repository.Repository) (chi.Router, error) {
	router := chi.NewRouter()
	authMiddleware := middleware.AuthMiddleware{Repo: repo}
	h := handlers.NewLoyaltyHandler(repo)

	router.Post(`/api/user/register`, h.UserRegister)
	router.Post(`/api/user/login`, h.UserLogin)
	router.Group(func(r chi.Router) {
		r.Use(authMiddleware.Login)
		r.Route(`/api/user`, func(r chi.Router) {
			r.Post(`/orders`, h.OrderUpload)
			r.Get(`/orders`, h.OrderList)
			r.Get(`/withdrawals`, h.BalanceWithdrawInfo)

			r.Route(`/balance`, func(r chi.Router) {
				r.Get(`/`, h.BalanceInfo)
				r.Post(`/withdraw`, h.BalanceWithdraw)
			})
		})
	})
	return router, nil
}
