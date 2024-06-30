package postgresql

import (
	"database/sql"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/accrual"
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/repository"
)

type RepoPostgreSQL struct {
	ConnString string
	accrual    accrual.BonusTracker
	db         *sql.DB
	lock       *sync.Mutex
}

func NewRepoPostgreSQL(conn string, accrual accrual.BonusTracker) (*RepoPostgreSQL, error) {
	storage := &RepoPostgreSQL{
		ConnString: conn,
		accrual:    accrual,
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
	createOrderTableSQL := `create table if not exists orders
(
    number      text                       not null
        constraint orders_pk
            primary key,
    user_id     integer                    not null
        constraint orders_users_id_fk
            references users
            on update cascade on delete cascade,
    status      text                       not null,
    accrual     double precision default 0 not null,
    uploaded_at text                       not null
);`
	createBalanceTableSQL := `create table if not exists balance
(
    user_id   integer                    not null
        constraint balance_pk
            primary key
        constraint balance_users_id_fk
            references users
            on update cascade on delete cascade,
    current   double precision default 0 not null,
    withdrawn double precision default 0 not null
);`
	createWithdrawTableSQL := `create table if not exists withdrawals
(
    id           serial
        constraint withdrawals_pk
            primary key,
    order_number text             not null,
    sum          double precision not null,
    processed_at text             not null,
    user_id      integer          not null
        constraint withdrawals_users_id_fk
            references users
            on update cascade on delete cascade
);`
	if _, err = db.Exec(createUserTableSQL); err != nil {
		defer func() { _ = db.Close() }()
		return err
	}
	if _, err = db.Exec(createOrderTableSQL); err != nil {
		defer func() { _ = db.Close() }()
		return err
	}
	if _, err = db.Exec(createBalanceTableSQL); err != nil {
		defer func() { _ = db.Close() }()
		return err
	}
	if _, err = db.Exec(createWithdrawTableSQL); err != nil {
		defer func() { _ = db.Close() }()
		return err
	}
	r.db = db
	return nil
}

var _ repository.Repository = &RepoPostgreSQL{}
