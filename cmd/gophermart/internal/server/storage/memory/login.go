package memory

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/repository"
)

func (r *RepoMemory) LoginUser(name, password string) (int, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	for _, u := range r.Users {
		if u.Name == name {
			if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
				return 0, repository.ErrInvalidPassword
			}
			return u.ID, nil
		}
	}
	return 0, repository.ErrUserNotFound
}

func (r *RepoMemory) CheckUserCredentials(name, password string) (int, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	for _, u := range r.Users {
		if u.Name == name {
			if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
				return 0, repository.ErrInvalidPassword
			}
			return u.ID, nil
		}
	}
	return 0, repository.ErrUserNotFound
}
