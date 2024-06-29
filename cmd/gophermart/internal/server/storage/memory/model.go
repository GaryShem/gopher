package memory

import (
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/logging"
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/repository"
)

type RepoMemory struct {
	Users              []repository.User
	UserIDToOrder      map[int][]repository.Order
	UserIDToBalance    map[int]repository.BalanceInfo
	UserIDToWithdrawal map[int][]repository.WithdrawalInfo
}

func NewRepoMemory(databaseURI string) *RepoMemory {
	logging.Log.Infof("database uri = %s, but memory storage ignores it", databaseURI)
	return &RepoMemory{
		Users:              []repository.User{},
		UserIDToOrder:      map[int][]repository.Order{},
		UserIDToBalance:    map[int]repository.BalanceInfo{},
		UserIDToWithdrawal: map[int][]repository.WithdrawalInfo{},
	}
}

var _ repository.Repository = &RepoMemory{}
