package models

import (
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/project/carecrew/config"
)

type GetReportInfo struct {
	Report_ID   int      `db:"report_id" json:"report_id"`
	Title       string   `db:"title" json:"title"`
	PersonnelID int      `db:"personnel_id" json:"personnel_id"`
	Detail      string   `db:"detail" json:"detail"`
	Location    string   `db:"location" json:"location"`
	Created_at  string   `db:"created_at" json:"created_at"`
	File        string   `db:"file" json:"-"`
	Files       []string `json:"files"`
}

func GetReport(db *sqlx.DB) ([]GetReportInfo, error) {
	gReport := []GetReportInfo{}

	query := `SELECT g.report_id, g.personnel_id, g.title, g.detail, g.location, g.file, g.created_at
	FROM "Reports" g
	ORDER BY g.report_id`

	err := db.Select(&gReport, query)
	if err != nil {
		return nil, err
	}

	for i := range gReport {
		if gReport[i].File != "" {
			files := strings.Split(gReport[i].File, ",")
			for j := range files {
				files[j] = config.APIURL + strings.TrimSpace(files[j])
			}
			gReport[i].Files = files
		}
	}

	return gReport, err
}
