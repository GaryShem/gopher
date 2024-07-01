package postgresql

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/logging"
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/repository"
)

func (r *RepoPostgreSQL) LoginUser(name, password string) (int, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	logging.Log.Infoln("User login attempt", name, password)
	user, err := r.GetUserByName(name)
	if err != nil {
		return 0, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return 0, repository.ErrInvalidPassword
	}
	return user.ID, nil
}

func (r *RepoPostgreSQL) CheckUserCredentials(name, password string) (int, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	user, err := r.GetUserByName(name)
	if err != nil {
		return 0, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return 0, repository.ErrInvalidPassword
	}
	return user.ID, nil
}
