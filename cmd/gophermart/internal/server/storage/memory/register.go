package memory

import (
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/logging"
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/repository"
)

func (r *RepoMemory) RegisterUser(name, password string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	for _, u := range r.Users {
		if u.Name == name {
			return repository.ErrUserAlreadyExists
		}
	}
	logging.Log.Infoln("Registering user", name, password)
	r.Users = append(r.Users, repository.User{
		ID:       len(r.Users) + 1,
		Name:     name,
		Password: password,
	})
	return nil
}
