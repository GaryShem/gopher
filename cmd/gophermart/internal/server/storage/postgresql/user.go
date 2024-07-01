package postgresql

import (
	"database/sql"
	"errors"

	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/repository"
)

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
