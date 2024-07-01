package postgresql

import (
	"context"
	"database/sql"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/accrual"
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/repository"
)

const createUserTableSQL = `create table if not exists users
(
    id       serial
        constraint users_pk
            primary key,
    name     text not null,
    password text not null
);`
const createOrderTableSQL = `create table if not exists orders
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
const createBalanceTableSQL = `create table if not exists balance
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
const createWithdrawTableSQL = `create table if not exists withdrawals
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

type RepoPostgreSQL struct {
	ConnString   string
	bonusTracker accrual.BonusTracker
	db           *sql.DB
	lock         *sync.RWMutex
}

func NewRepoPostgreSQL(conn string, accrual accrual.BonusTracker) (*RepoPostgreSQL, error) {
	storage := &RepoPostgreSQL{
		ConnString:   conn,
		bonusTracker: accrual,
		lock:         &sync.RWMutex{},
	}
	err := storage.Init()
	if err != nil {
		return nil, err
	}
	return storage, nil
}

func (r *RepoPostgreSQL) createUserTable(tx *sql.Tx) error {
	_, err := tx.Exec(createUserTableSQL)
	return err
}

func (r *RepoPostgreSQL) createOrderTable(tx *sql.Tx) error {
	_, err := tx.Exec(createOrderTableSQL)
	return err
}

func (r *RepoPostgreSQL) createBalanceTable(tx *sql.Tx) error {
	_, err := tx.Exec(createBalanceTableSQL)
	return err
}

func (r *RepoPostgreSQL) createWithdrawTable(tx *sql.Tx) error {
	_, err := tx.Exec(createWithdrawTableSQL)
	return err
}

func (r *RepoPostgreSQL) Init() error {
	r.lock.Lock()
	defer r.lock.Unlock()
	db, err := sql.Open("pgx", r.ConnString)
	if err != nil {
		return err
	}
	r.db = db
	tx, err := r.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	// create tables
	if err = r.createUserTable(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err = r.createOrderTable(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err = r.createBalanceTable(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err = r.createWithdrawTable(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

var _ repository.Repository = &RepoPostgreSQL{}
