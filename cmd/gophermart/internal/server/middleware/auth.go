package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/logging"
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/repository"
)

const UserIDHeader = "used-id"

type AuthMiddleware struct {
	Repo repository.Repository
}

func (am *AuthMiddleware) Login(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startRequest := time.Now()
		username, password, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		id, err := am.Repo.CheckUserCredentials(username, password)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		r.Header.Set(UserIDHeader, fmt.Sprintf("%v", id))
		next.ServeHTTP(w, r)
		endRequest := time.Now()
		duration := endRequest.Sub(startRequest)
		logging.Log.Infoln("request took", duration)
	})
}
