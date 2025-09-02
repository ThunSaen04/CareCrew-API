package models

import (
	"github.com/jmoiron/sqlx"
)

type GuestReportInfo struct {
	Report_ID  int    `db:"report_id" json:"report_id"`   //Reports
	Title      string `db:"title" json:"title"`           //Reports
	Email      string `db:"email" json:"email"`           //Reports
	Detail     string `db:"detail" json:"detail"`         //Reports
	Location   string `db:"location" json:"location"`     //Reports
	Created_at string `db:"created_at" json:"created_at"` //Reports
	File       string `db:"file" json:"file"`             //Reports
}

func GuestReport(db *sqlx.DB, greportInfo *GuestReportInfo) error {

	_, err := db.Exec(
		`INSERT INTO "Reports" (title, email, detail, location, file, created_at)
		VALUES ($1, $2, $3, $4, $5,NOW())`,
		greportInfo.Title,
		greportInfo.Email,
		greportInfo.Detail,
		greportInfo.Location,
		greportInfo.File,
	)
	if err != nil {
		return err
	}

	return nil
}
