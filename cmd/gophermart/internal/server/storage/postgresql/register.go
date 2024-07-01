package postgresql

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/logging"
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/repository"
)

func (r *RepoPostgreSQL) RegisterUser(name, password string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	_, err := r.GetUserByName(name)
	if !errors.Is(err, repository.ErrUserNotFound) {
		return err
	}
	logging.Log.Infoln("Registering user", name, password)
	userQueryTemplate := `INSERT INTO users 
    	(name, password) 
		VALUES (@name, @password)
		RETURNING id`
	args := pgx.NamedArgs{
		"name":     name,
		"password": password,
	}
	var userID int
	err = r.db.QueryRow(userQueryTemplate, args).Scan(&userID)
	if err != nil {
		return err
	}

	logging.Log.Infoln("Creating balance for user", name)
	balanceQueryTemplate := `INSERT INTO balance (user_id) VALUES (@user_id)`
	args = pgx.NamedArgs{
		"user_id": userID,
	}
	_, err = r.db.Exec(balanceQueryTemplate, args)
	if err != nil {
		return err
	}
	return nil
}
