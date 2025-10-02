package models

import (
	"encoding/json"
	"errors"

	"github.com/jmoiron/sqlx"
)

type NotificationInfo struct {
	Noti_id      int             `db:"noti_id" json:"noti_id"`
	Personnel_Id int             `db:"personnel_id" json:"personnel_id"`
	Title        string          `db:"title" json:"title"`
	Body         string          `db:"body" json:"body"`
	Read         bool            `db:"read" json:"read"`
	Data         json.RawMessage `db:"data" json:"data" swaggertype:"object"`
	Created_at   string          `db:"created_at" json:"created_at"`
}

func GetNoti(db *sqlx.DB) ([]NotificationInfo, error) {
	notis := []NotificationInfo{}

	query := `
		SELECT noti_id, personnel_id, title, body, data, read, created_at
		FROM "Notifications"
		ORDER BY created_at DESC
	`
	err := db.Select(&notis, query)
	if err != nil {
		return nil, err
	}

	return notis, nil
}

type ReadNotisInfo struct {
	Noti_id      int `db:"noti_id" json:"noti_id"`
	Personnel_Id int `db:"personnel_id" json:"personnel_id"`
}

func ReadNotis(db *sqlx.DB, readnotisinfo ReadNotisInfo) error {
	tranX, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tranX.Rollback()

	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM "Notifications"
			WHERE noti_id = $1 AND personnel_id = $2
		)
	`
	err = tranX.Get(&exists, query, readnotisinfo.Noti_id, readnotisinfo.Personnel_Id)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("ไม่พบประวัติแจ้งเตือน")
	}

	_, err = tranX.Exec(`
		UPDATE "Notifications"
		SET read = TRUE
		WHERE noti_id = $1 AND personnel_id = $2
	`, readnotisinfo.Noti_id, readnotisinfo.Personnel_Id)
	if err != nil {
		return err
	}

	return tranX.Commit()
}
