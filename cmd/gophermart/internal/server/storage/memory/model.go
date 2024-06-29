package memory

import (
	"sync"

	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/logging"
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/repository"
)

type RepoMemory struct {
	Users              []repository.User
	UserIDToOrder      map[int][]repository.Order
	UserIDToBalance    map[int]repository.BalanceInfo
	UserIDToWithdrawal map[int][]repository.WithdrawalInfo
	lock               sync.Mutex
}

func NewRepoMemory(databaseURI string) (*RepoMemory, error) {
	logging.Log.Infof("database uri = %s, but memory storage ignores it", databaseURI)
	return &RepoMemory{
		Users:              []repository.User{},
		UserIDToOrder:      map[int][]repository.Order{},
		UserIDToBalance:    map[int]repository.BalanceInfo{},
		UserIDToWithdrawal: map[int][]repository.WithdrawalInfo{},
	}, nil
}

var _ repository.Repository = &RepoMemory{}
