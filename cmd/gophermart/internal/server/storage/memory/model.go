package memory

import (
	"sync"

	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/accrual"
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/logging"
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/repository"
)

type RepoMemory struct {
	Users              []repository.User
	UserIDToOrder      map[int]map[string]repository.Order
	UserIDToBalance    map[int]repository.BalanceInfo
	UserIDToWithdrawal map[int][]repository.WithdrawalInfo
	accrual            accrual.BonusTracker
	lock               sync.Mutex
}

func NewRepoMemory(databaseURI string, accrual accrual.BonusTracker) (*RepoMemory, error) {
	logging.Log.Infof("database uri = %s, but memory storage ignores it", databaseURI)
	return &RepoMemory{
		Users:              []repository.User{},
		UserIDToOrder:      map[int]map[string]repository.Order{},
		UserIDToBalance:    map[int]repository.BalanceInfo{},
		UserIDToWithdrawal: map[int][]repository.WithdrawalInfo{},
		accrual:            accrual,
	}, nil
}

var _ repository.Repository = &RepoMemory{}
