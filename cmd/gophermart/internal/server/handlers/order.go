package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/repository"
)

func (l *LoyaltyHandler) OrderUpload(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.Header.Get("user_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	orderNumberBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	orderNumber := string(orderNumberBytes)
	err = l.repo.OrderUpload(userID, orderNumber)
	if err != nil {
		if errors.Is(err, repository.ErrOrderUploadedSameUser) {
			w.WriteHeader(http.StatusOK)
		} else if errors.Is(err, repository.ErrOrderUploadedDifferentUser) {
			w.WriteHeader(http.StatusConflict)
		} else if errors.Is(err, repository.ErrOrderIDFormatInvalid) {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (l *LoyaltyHandler) OrderList(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	orders, err := l.repo.OrderGet(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
	}
	jsonData, err := json.Marshal(orders)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(jsonData)
	return
}
