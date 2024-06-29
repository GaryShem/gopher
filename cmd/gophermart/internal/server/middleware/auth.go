package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/repository"
)

const USER_ID_HEADER = "used-id"

type AuthMiddleware struct {
	Repo repository.Repository
}

func (am *AuthMiddleware) Login(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		username, password, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := am.Repo.UserLogin(username, password)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		r.Header.Set(USER_ID_HEADER, fmt.Sprintf("%v", id))
		next.ServeHTTP(w, r)
	})
}
