package handlers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/repository"
)

func (l *LoyaltyHandler) UserRegister(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer func() { _ = r.Body.Close() }()
	var register repository.RegisterRequest
	if err = json.Unmarshal(body, &register); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = l.repo.UserRegister(register.Login, register.Password); err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	b64token := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", register.Login, register.Password)))
	w.Header().Set("Authorization", fmt.Sprintf("Basic %s", b64token))
	w.WriteHeader(http.StatusOK)
}

func (l *LoyaltyHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer func() { _ = r.Body.Close() }()
	var register repository.RegisterRequest
	if err = json.Unmarshal(body, &register); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err = l.repo.UserLogin(register.Login, register.Password); err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	b64token := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", register.Login, register.Password)))
	w.Header().Set("Authorization", fmt.Sprintf("Basic %s", b64token))
	w.WriteHeader(http.StatusOK)
}
