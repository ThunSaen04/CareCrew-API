package orther

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type UpdateProfileInfo struct {
	PersonnelID int    `db:"personnel_id" json:"personnel_id"`
	File        string `db:"file" json:"file"`
}

func UpdateProfile(db *sqlx.DB, updateprofileinfo UpdateProfileInfo) error {

	tranX, err := db.Beginx()
	if err != nil {
		return err
	}

	defer tranX.Rollback()

	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1 FROM "Personnels"
			WHERE personnel_id = $1
		)
	`
	err = tranX.Get(&exists, query, updateprofileinfo.PersonnelID)
	if err != nil {
		return err
	}

	if exists {
		query = `
			UPDATE "Personnels" 
			SET file = $1 
			WHERE personnel_id = $2
		`
		_, err = tranX.Exec(query, updateprofileinfo.File, updateprofileinfo.PersonnelID)
		if err != nil {
			return err
		}
	} else {
		return errors.New("ไม่พบหมายเลขผู้ใช้งานดังกล่าว")
	}

	err = tranX.Commit()
	if err != nil {
		return err
	}

	return nil
}
