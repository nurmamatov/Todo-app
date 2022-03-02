package storage

import (
	"github.com/jmoiron/sqlx"
	"app/storage/postgres"
	"app/storage/repo"
)

// I is an interface for storage
type I interface {
	SendEmail() repo.SendStorageI
}

type storagePg struct {
	db        *sqlx.DB
	sendRepo repo.SendStorageI
}

// NewStoragePg ...
func NewStoragePg(db *sqlx.DB) I {
	return &storagePg{
		db:        db,
		sendRepo: postgres.NewSendRepo(db),
	}
}

func (s storagePg) SendEmail() repo.SendStorageI {
	return s.sendRepo
}
