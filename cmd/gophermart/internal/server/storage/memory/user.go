package memory

import (
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/logging"
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/repository"
)

func (r *RepoMemory) UserRegister(name, password string) error {
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

func (r *RepoMemory) UserLogin(name, password string) (int, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	logging.Log.Infoln("User login attempt", name, password)
	for _, u := range r.Users {
		if u.Name == name {
			if u.Password == password {
				return u.ID, nil
			}
			return 0, repository.ErrInvalidPassword
		}
	}
	return 0, repository.ErrUserNotFound
}

func (r *RepoMemory) GetUserByName(name string) (repository.User, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	for _, u := range r.Users {
		if u.Name == name {
			return u, nil
		}
	}
	return repository.User{}, repository.ErrUserNotFound
}
