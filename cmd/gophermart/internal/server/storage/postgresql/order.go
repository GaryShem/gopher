package postgresql

import (
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/repository"
)

func (r *RepoPostgreSQL) OrderUpload(userID int, orderID string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	panic("implement me")
}

func (r *RepoPostgreSQL) GetOrdersByUserID(userID int) ([]repository.Order, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	panic("implement me")
}

func (r RepoPostgreSQL) UpdateOrderProcessing(orderID string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	panic("implement me")
}
