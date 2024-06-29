package memory

import (
	"time"

	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/repository"
)

func (r *RepoMemory) OrderUpload(userID int, orderID string) error {
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
		orders = []repository.Order{}
	}
	r.UserIDToOrder[userID] = append(orders, repository.Order{
		Number:     orderID,
		Status:     "NEW",
		UploadedAt: time.Now().Format(time.RFC3339),
	})
	return nil
}

func (r *RepoMemory) OrderGet(userID int) ([]repository.Order, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	orders, ok := r.UserIDToOrder[userID]
	if !ok {
		return []repository.Order{}, nil
	}
	return orders, nil
}

func (r *RepoMemory) GetOrdersByUserID(userID int) ([]repository.Order, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	orders, ok := r.UserIDToOrder[userID]
	if !ok {
		return []repository.Order{}, nil
	}
	return orders, nil
}
