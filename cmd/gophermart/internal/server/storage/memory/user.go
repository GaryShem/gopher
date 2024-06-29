package memory

import "github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/repository"

func (r *RepoMemory) UserRegister(name, password string) error {
	for _, u := range r.Users {
		if u.Name == name {
			return repository.ErrUserAlreadyExists
		}
	}
	r.Users = append(r.Users, repository.User{
		ID:       len(r.Users) + 1,
		Name:     name,
		Password: password,
	})
	return nil
}

func (r *RepoMemory) UserLogin(name, password string) (int, error) {
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
