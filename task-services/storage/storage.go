package storage

import (
	"two_services/task-services/storage/postgres"
	"two_services/task-services/storage/repo"

	"github.com/jmoiron/sqlx"
)

type IStorage interface {
	Task() repo.TaskStorageI
	Email() repo.EmailStorageI
}

type storagePg struct {
	db        *sqlx.DB
	taskRepo  repo.TaskStorageI
	emailRepo repo.EmailStorageI
}

func NewStoragePg(db *sqlx.DB) *storagePg {
	return &storagePg{
		db:       db,
		taskRepo: postgres.NewTaskRepo(db),
	}
}

func (s storagePg) Task() repo.TaskStorageI {
	return s.taskRepo
}
func (s storagePg) Email() repo.EmailStorageI {
	return s.emailRepo
}
