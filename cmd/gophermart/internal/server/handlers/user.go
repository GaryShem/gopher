package handlers

import (
	"errors"
	"net/http"

	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/repository"
)

func (l *LoyaltyHandler) UserRegister(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := l.repo.UserRegister(username, password); err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (l *LoyaltyHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err := l.repo.UserLogin(username, password)
	if errors.Is(err, repository.ErrInvalidPassword) || errors.Is(err, repository.ErrUserNotFound) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
