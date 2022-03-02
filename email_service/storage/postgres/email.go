package postgres

import (
	"app/storage/repo"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type sendRepo struct {
	db *sqlx.DB
}

// NewSendRepo ...
func NewSendRepo(db *sqlx.DB) repo.SendStorageI {
	return &sendRepo{db: db}
}

func (cm *sendRepo) MakeSent(ID string) error {
	var err error
	makesent := `UPDATE email_send_email SET send_status=true where id = $1`
	cm.db.MustExec(makesent, ID)
	return err
}

func (cm *sendRepo) Send(subject, body string, status bool, val string) error {
	textID, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	insert := `
	INSERT INTO
	email_text
	(
		id,
		subject,
		body,
		status

	)
	values($1, $2, $3,$4)
	`
	_, err = cm.db.Exec(insert, textID, subject, body, status)
	if err != nil {
		return err
	}

	return nil
}
func (cm *sendRepo) SendSms(phone_number, text string) error {
	textID, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	insert := `
	INSERT INTO
	sms_text
	(
		id,
		text,
		phone_number
	)
	values($1, $2, $3)
	`
	_, err = cm.db.Exec(insert, textID, phone_number,text)
	if err != nil {
		return err
	}

	return nil
}
