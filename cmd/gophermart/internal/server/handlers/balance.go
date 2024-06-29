package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (l *LoyaltyHandler) BalanceInfo(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.Header.Get("user_id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	info, err := l.repo.BalanceList(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	json, err := json.Marshal(info)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(json)
	return
}

func (l *LoyaltyHandler) BalanceWithdraw(w http.ResponseWriter, r *http.Request) {
	panic("Not implemented")
}

func (l *LoyaltyHandler) BalanceWithdrawInfo(w http.ResponseWriter, r *http.Request) {
	panic("Not implemented")
}
