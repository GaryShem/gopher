package memory

import (
	"fmt"
	"time"

	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/repository"
)

func (r *RepoMemory) ListBalance(userID int) (repository.BalanceInfo, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	balance, ok := r.UserIDToBalance[userID]
	if !ok {
		return repository.BalanceInfo{}, repository.ErrUserNotFound
	}
	return balance, nil
}
func (r *RepoMemory) WithdrawBalance(userID int, orderID string, amount float64) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	if err := repository.ValidateOrderID(orderID); err != nil {
		return repository.ErrOrderIDFormatInvalid
	}
	balance, ok := r.UserIDToBalance[userID]
	if !ok {
		return repository.ErrUserNotFound
	}
	if balance.Current < amount {
		return fmt.Errorf("%w withdrawing %v from %v", repository.ErrBalanceNotEnough, amount, balance.Current)
	}
	balance.Current -= amount
	balance.Withdrawn += amount
	r.UserIDToBalance[userID] = balance
	r.UserIDToWithdrawal[userID] = append(r.UserIDToWithdrawal[userID], repository.WithdrawalInfo{
		Order:       orderID,
		Sum:         amount,
		ProcessedAt: time.Now().UTC().Format(time.RFC3339),
	})
	return nil
}
func (r *RepoMemory) GetBalanceWithdrawInfo(userID int) ([]repository.WithdrawalInfo, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	info, ok := r.UserIDToWithdrawal[userID]
	if !ok {
		return []repository.WithdrawalInfo{}, repository.ErrNoWithdrawals
	}
	return info, nil
}
