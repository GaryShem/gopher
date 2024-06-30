package postgresql

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/logging"
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/repository"
)

func (r *RepoPostgreSQL) UserRegister(name, password string) error {
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

func (r *RepoPostgreSQL) UserLogin(name, password string) (int, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	logging.Log.Infoln("User login attempt", name, password)
	user, err := r.GetUserByName(name)
	if err != nil {
		return 0, err
	}
	if user.Password != password {
		return 0, repository.ErrInvalidPassword
	}
	return user.ID, nil
}

func (r *RepoPostgreSQL) GetUserByName(name string) (repository.User, error) {
	queryTemplate := `SELECT * FROM users WHERE name = $1`

	res := r.db.QueryRow(queryTemplate, name)
	var user repository.User
	if err := res.Scan(&user.ID, &user.Name, &user.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.User{}, repository.ErrUserNotFound
		} else {
			return repository.User{}, err
		}
	}
	return user, nil
}
