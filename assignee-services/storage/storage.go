package storage

import (
	"two_services/assignee-services/storage/postgres"
	"two_services/assignee-services/storage/repo"

	"github.com/jmoiron/sqlx"
)

type IStorage interface {
	User() repo.UserStorageI
}

type storagePg struct {
	db       *sqlx.DB
	taskRepo repo.UserStorageI
}

func NewStoragePg(db *sqlx.DB) *storagePg {
	return &storagePg{
		db:       db,
		taskRepo: postgres.NewUserRepo(db),
	}
}

func (s storagePg) User() repo.UserStorageI {
	return s.taskRepo
}
