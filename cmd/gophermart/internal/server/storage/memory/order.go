package memory

import (
	"errors"
	"sort"
	"time"

	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/logging"
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/repository"
)

var ErrInternalConsistencyError = errors.New("internal consistency error")

func (r *RepoMemory) UploadOrder(userID int, orderID string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	if err := repository.ValidateOrderID(orderID); err != nil {
		return err
	}
	for u := range r.UserIDToOrder {
		for _, o := range r.UserIDToOrder[u] {
			if o.Number == orderID {
				if u == userID {
					return repository.ErrOrderUploadedSameUser
				}
				return repository.ErrOrderUploadedDifferentUser
			}
		}
	}
	orders, ok := r.UserIDToOrder[userID]
	if !ok {
		orders = map[string]repository.Order{}
	}
	orders[orderID] = repository.Order{
		Number:     orderID,
		Status:     "NEW",
		UploadedAt: time.Now().UTC().Format(time.RFC3339),
	}
	r.UserIDToOrder[userID] = orders
	go r.ProcessOrderUpdate(userID, orderID)
	return nil
}

func (r *RepoMemory) GetOrdersByUser(userID int) ([]repository.Order, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	orders, ok := r.UserIDToOrder[userID]
	result := make([]repository.Order, 0)
	if ok {
		for _, order := range orders {
			result = append(result, order)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].UploadedAt < result[j].UploadedAt
	})
	return result, nil
}

func (r *RepoMemory) ProcessOrderUpdate(userID int, orderID string) error {
	startUpdate := time.Now()
	defer func() {
		logging.Log.Infoln("ProcessOrderUpdate took", time.Since(startUpdate))
	}()
	info, err := r.bonusTracker.UpdateOrder(orderID)
	if err != nil {
		return err
	}

	r.lock.Lock()
	defer r.lock.Unlock()
	r.UserIDToOrder[userID][orderID] = *info
	for u := range r.UserIDToOrder {
		_, ok := r.UserIDToOrder[u][orderID]
		if ok {
			r.UserIDToOrder[u][orderID] = *info
			balance := r.UserIDToBalance[u]
			balance.Current += info.Accrual
			r.UserIDToBalance[u] = balance
			logging.Log.Infoln("ProcessOrderUpdate success")
			return nil
		}
	}
	logging.Log.Infoln("ProcessOrderUpdate fail")
	return ErrInternalConsistencyError
}
