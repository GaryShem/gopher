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

func (l *LoyaltyHandler) OrderUpload(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.Header.Get(middleware.UserIDHeader))
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
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (l *LoyaltyHandler) OrderList(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.Header.Get(middleware.UserIDHeader))
	if err != nil {
		w.WriteHeader(503)
		return
	}
	orders, err := l.repo.GetOrdersByUserID(userID)
	if err != nil {
		w.WriteHeader(502)
		return
	}
	//if len(orders) == 0 {
	//	w.Header().Set("Content-Type", "application/json")
	//	w.WriteHeader(http.StatusNoContent)
	//	return
	//}
	jsonData, err := json.Marshal(orders)
	if err != nil {
		w.WriteHeader(501)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	_, _ = w.Write(jsonData)
}
