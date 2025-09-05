package models

import (
	"log"

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

func GuestReport(db *sqlx.DB, greportInfo *ReportInfo) error {

	log.Println(greportInfo.Title)

	_, err := db.Exec(
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

	return nil
}
