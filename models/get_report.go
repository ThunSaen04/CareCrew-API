package models

import "github.com/jmoiron/sqlx"

func GetReport(db *sqlx.DB) ([]GuestReportInfo, error) {
	gReport := []GuestReportInfo{}

	query := `SELECT g.report_id, g.email, g.detail, g.location, g.created_at
	FROM "Reports" g
	ORDER BY g.report_id`

	err := db.Select(&gReport, query)

	return gReport, err
}
