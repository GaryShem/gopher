package memory

import (
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/repository"
)

func (r *RepoMemory) GetUserByName(name string) (repository.User, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	for _, u := range r.Users {
		if u.Name == name {
			return u, nil
		}
	}
	return repository.User{}, repository.ErrUserNotFound
}
