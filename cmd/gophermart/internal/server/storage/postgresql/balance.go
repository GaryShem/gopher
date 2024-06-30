package postgresql

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/repository"
)

func (r *RepoPostgreSQL) BalanceList(userID int) (repository.BalanceInfo, error) {
	selectQuery := `SELECT current, withdrawn FROM balance WHERE user_id = @user_id`
	args := pgx.NamedArgs{
		"user_id": userID,
	}
	var info repository.BalanceInfo
	row := r.db.QueryRow(selectQuery, args)
	err := row.Scan(&info.Current, &info.Withdrawn)
	return info, err
}
func (r *RepoPostgreSQL) BalanceWithdraw(userID int, orderID string, amount float64) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	if err := repository.ValidateOrderID(orderID); err != nil {
		return repository.ErrOrderIDFormatInvalid
	}
	balance, err := r.BalanceList(userID)
	if err != nil {
		return err
	}
	if balance.Current < amount {
		return fmt.Errorf("%w withdrawing %v from %v", repository.ErrBalanceNotEnough, amount, balance.Current)
	}
	balance.Current -= amount
	balance.Withdrawn += amount
	balanceUpdateQuery := `UPDATE balance SET
		current = @current, withdrawn = @withdrawn
        WHERE user_id = @user_id`
	args := pgx.NamedArgs{
		"user_id":   userID,
		"current":   balance.Current,
		"withdrawn": balance.Withdrawn,
	}

	if _, err = r.db.Exec(balanceUpdateQuery, args); err != nil {
		return err
	}
	withdrawInsertQuery := `INSERT INTO withdrawals 
    	(order_number, sum, processed_at, user_id) VALUES 
    	(@order_number, @sum, @processed_at, @user_id)`
	args = pgx.NamedArgs{
		"order_number": orderID,
		"user_id":      userID,
		"processed_at": time.Now().Format(time.RFC3339),
		"sum":          amount,
	}
	if _, err = r.db.Exec(withdrawInsertQuery, args); err != nil {
		return err
	}
	return nil
}
func (r *RepoPostgreSQL) BalanceWithdrawInfo(userID int) ([]repository.WithdrawalInfo, error) {
	selectQuery := `SELECT order_number, sum, processed_at FROM withdrawals WHERE user_id = @user_id`
	args := pgx.NamedArgs{
		"user_id": userID,
	}
	rows, err := r.db.Query(selectQuery, args)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	result := []repository.WithdrawalInfo{}
	for rows.Next() {
		var info repository.WithdrawalInfo
		err = rows.Scan(&info.Order, &info.Sum, &info.ProcessedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, info)
	}

	return result, err
}
