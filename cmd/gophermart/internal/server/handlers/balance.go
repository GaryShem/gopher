package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/middleware"
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/repository"
)

func (l *LoyaltyHandler) BalanceInfo(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.Header.Get(middleware.UserIdHeader))
	if err != nil {
		w.WriteHeader(501)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	info, err := l.repo.BalanceList(userID)
	if err != nil {
		w.WriteHeader(502)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	json, err := json.Marshal(info)
	if err != nil {
		w.WriteHeader(503)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(json)
	return
}

func (l *LoyaltyHandler) BalanceWithdraw(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.Header.Get(middleware.UserIdHeader))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("unable to read request body"))
		return
	}
	info := new(repository.WithdrawalInfo)
	err = json.Unmarshal(body, info)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("unable to unmarshal request body"))
		return
	}
	err = l.repo.BalanceWithdraw(userID, info.Order, info.Sum)
	if err != nil {
		if errors.Is(err, repository.ErrBalanceNotEnough) {
			w.WriteHeader(http.StatusPaymentRequired)
		} else if errors.Is(err, repository.ErrOrderIDFormatInvalid) {
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (l *LoyaltyHandler) BalanceWithdrawInfo(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.Header.Get(middleware.UserIdHeader))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	withdrawals, err := l.repo.BalanceWithdrawInfo(userID)
	if err != nil {
		if errors.Is(err, repository.ErrNoWithdrawals) {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	jsonObj, err := json.Marshal(withdrawals)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("unable to marshal response body"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonObj)
	w.WriteHeader(http.StatusOK)
}
