package postgresql

import (
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/repository"
)

func (r *RepoPostgreSQL) BalanceList(userID int) (repository.BalanceInfo, error) {
	panic("implement me")
}
func (r *RepoPostgreSQL) BalanceWithdraw(userID int, orderID string, amount float64) error {
	panic("implement me")
}
func (r *RepoPostgreSQL) BalanceWithdrawInfo(userID int) ([]repository.WithdrawalInfo, error) {
	panic("implement me")
}
