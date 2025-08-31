package models

import (
	"github.com/jmoiron/sqlx"
)

type GuestReportInfo struct {
	Report_ID  int    `db:"report_id" json:"report_id"`   //Guest_Reports
	Title      string `db:"title" json:"title"`           //Guest_Reports
	Email      string `db:"email" json:"email"`           //Guest_Reports
	Detail     string `db:"detail" json:"detail"`         //Guest_Reports
	Location   string `db:"location" json:"location"`     //Guest_Reports
	Created_at string `db:"created_at" json:"created_at"` //Guest_Reports
	File       string `db:"file" json:"file"`             //Guest_Reports
}

func GuestReport(db *sqlx.DB, greportInfo *GuestReportInfo) error {

	_, err := db.Exec(
		`INSERT INTO "Guest_Reports" (title, email, detail, location, file, created_at)
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
