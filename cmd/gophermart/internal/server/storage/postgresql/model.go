package postgresql

import (
	"database/sql"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/repository"
)

type RepoPostgreSQL struct {
	ConnString string
	db         *sql.DB
	lock       *sync.Mutex
}

func NewRepoPostgreSQL(conn string) (*RepoPostgreSQL, error) {
	storage := &RepoPostgreSQL{
		ConnString: conn,
		lock:       &sync.Mutex{},
	}
	err := storage.Init()
	if err != nil {
		return nil, err
	}
	return storage, nil
}

func (r *RepoPostgreSQL) Init() error {
	db, err := sql.Open("pgx", r.ConnString)
	if err != nil {
		return err
	}
	// create tables
	createUserTableSQL := `create table if not exists users
(
    id       serial
        constraint users_pk
            primary key,
    name     text not null,
    password text not null
);`
	if _, err = db.Exec(createUserTableSQL); err != nil {
		defer func() { _ = db.Close() }()
		return err
	}
	r.db = db
	return nil
}

var _ repository.Repository = &RepoPostgreSQL{}
