package models

import (
	"github.com/jmoiron/sqlx"
)

type ReportInfo struct {
	Report_ID   int    `db:"report_id" json:"report_id"`       //Reports
	Title       string `db:"title" json:"title"`               //Reports
	PersonnelID int    `db:"personnel_id" json:"personnel_id"` //Reports
	Detail      string `db:"detail" json:"detail"`             //Reports
	Location    string `db:"location" json:"location"`         //Reports
	Created_at  string `db:"created_at" json:"created_at"`     //Reports
	File        string `db:"file" json:"file"`                 //Reports
}

func Report(db *sqlx.DB, greportInfo *ReportInfo) error {

	tranX, err := db.Beginx()
	if err != nil {
		return err
	}

	defer tranX.Rollback()

	_, err = tranX.Exec(
		`INSERT INTO "Reports" (title, personnel_id, detail, location, file, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())`,
		greportInfo.Title,
		greportInfo.PersonnelID,
		greportInfo.Detail,
		greportInfo.Location,
		greportInfo.File,
	)
	if err != nil {
		return err
	}

	_, err = tranX.Exec(`
			INSERT INTO "Worker_logs" (personnel_id, report_id, detail, file)
			VALUES ($1, $2, $3, $4)
		`, greportInfo.PersonnelID, greportInfo.Report_ID, "แจ้งเหตุสร้างงาน", greportInfo.File)
	if err != nil {
		return err
	}

	err = tranX.Commit()
	if err != nil {
		return err
	}

	return nil
}
