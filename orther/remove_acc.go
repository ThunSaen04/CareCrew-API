package orther

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type RemoveAccInfo struct {
	PersonnelID int `db:"personnel_id" json:"personnel_id"`
}

func RemoveAcc(db *sqlx.DB, removeAccInfo *RemoveAccInfo) error {
	tranX, err := db.Beginx()
	if err != nil {
		return err
	}

	defer tranX.Rollback()

	var exists bool
	//เช็คว่ามีบัญชีนั้นจริงอ๊ะป่าว
	query := `
		SELECT EXISTS (
			SELECT 1 FROM "Personnels"
			WHERE personnel_id = $1
		)
	`
	err = tranX.Get(&exists, query, removeAccInfo.PersonnelID)
	if err != nil {
		return err
	}

	if exists {
		_, err = tranX.Exec(`
			DELETE FROM "Personnels"
			WHERE personnel_id = $1
		`, removeAccInfo.PersonnelID)
		if err != nil {
			return err
		}
	} else {
		return errors.New("ไม่พบบัญชีผู้ใช้")
	}

	err = tranX.Commit()
	if err != nil {
		return err
	}
	return nil
}
