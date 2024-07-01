package postgresql

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/logging"
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/repository"
)

func (r *RepoPostgreSQL) UploadOrder(userID int, orderID string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	if err := repository.ValidateOrderID(orderID); err != nil {
		return err
	}
	var existingOrder repository.Order
	orderSelectSQL := `SELECT user_id FROM orders WHERE number = @order_id`
	selectArgs := pgx.NamedArgs{
		"order_id": orderID,
	}
	err := r.db.QueryRow(orderSelectSQL, selectArgs).Scan(&existingOrder.UserID)
	if err == nil {
		if existingOrder.UserID == userID {
			return repository.ErrOrderUploadedSameUser
		} else {
			return repository.ErrOrderUploadedDifferentUser
		}
	} else {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}
	orderInsertSQL := `INSERT INTO orders 
    	(number, user_id, status, accrual, uploaded_at)
		VALUES (@order_id, @user_id, @status, @accrual, @upload_time)
		RETURNING number`
	insertArgs := pgx.NamedArgs{
		"order_id":    orderID,
		"user_id":     userID,
		"status":      "NEW",
		"accrual":     0,
		"upload_time": time.Now().UTC().Format(time.RFC3339),
	}
	var addedOrder repository.Order
	if err = r.db.QueryRow(orderInsertSQL, insertArgs).Scan(&addedOrder.Number); err != nil {
		return err
	}
	go r.ProcessOrderUpdate(userID, orderID)
	return nil
}

func (r *RepoPostgreSQL) GetOrdersByUser(userID int) ([]repository.Order, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	orderSelectSQL := `SELECT number, status, accrual, uploaded_at FROM orders 
         WHERE user_id = @user_id
         ORDER BY uploaded_at`
	selectArgs := pgx.NamedArgs{
		"user_id": userID,
	}
	result := make([]repository.Order, 0)
	rows, err := r.db.Query(orderSelectSQL, selectArgs)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, nil
		} else {
			return nil, err
		}
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var order repository.Order
		if err = rows.Scan(
			&order.Number,
			&order.Status,
			&order.Accrual,
			&order.UploadedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, order)
	}
	return result, nil
}

func (r *RepoPostgreSQL) ProcessOrderUpdate(userID int, orderID string) error {
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
	orderUpdateQuery := `UPDATE orders SET
    	status = @status, accrual = @accrual 
    	WHERE number = @order_id`
	args := pgx.NamedArgs{
		"order_id": orderID,
		"status":   info.Status,
		"accrual":  info.Accrual,
	}
	if _, err = r.db.Exec(orderUpdateQuery, args); err != nil {
		return err
	}
	if info.Accrual > 0 {
		args = pgx.NamedArgs{
			"user_id": userID,
		}
		var balance repository.BalanceInfo
		row := r.db.QueryRow(selectBalanceQuery, args)
		err = row.Scan(&balance.Current, &balance.Withdrawn)
		if err != nil {
			return err
		}
		newBalance := balance.Current + info.Accrual
		balanceUpdateQuery := `UPDATE balance SET
			current = @current 
            WHERE user_id = @user_id`
		args = pgx.NamedArgs{
			"user_id": userID,
			"current": newBalance,
		}

		if _, err = r.db.Exec(balanceUpdateQuery, args); err != nil {
			return err
		}
	}
	return nil
}
