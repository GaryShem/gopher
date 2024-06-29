package router

import (
	"net/http"

	chi "github.com/go-chi/chi/v5"

	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/handlers"
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/middleware"
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/repository"
)

func GopherRouter(repo repository.Repository, middlewares ...func(http.Handler) http.Handler) (chi.Router, error) {
	router := chi.NewRouter()
	authMiddleware := middleware.AuthMiddleware{repo}
	h := handlers.NewLoyaltyHandler(repo)
	for _, mw := range middlewares {
		router.Use(mw)
	}
	router.Post(`/api/user/register`, h.UserRegister)
	router.Group(func(r chi.Router) {
		r.Use(authMiddleware.Login)
		r.Route(`/api/user`, func(r chi.Router) {
			r.Post(`/register`, h.UserRegister)
			r.Post(`/login`, h.UserLogin)

			r.Post(`/orders`, h.OrderUpload)
			r.Get(`/orders`, h.OrderList)

			r.Route(`/balance`, func(r chi.Router) {
				r.Get(`/`, h.BalanceInfo)
				r.Post(`/withdraw`, h.BalanceWithdraw)
				r.Get(`/withdrawals`, h.BalanceWithdrawInfo)
			})
		})
	})
	return router, nil
}
