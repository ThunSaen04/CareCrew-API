package assignor

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type RemoveReportInfo struct {
	Report_id int `db:"report_id" json:"report_id"`
}

func RemoveReport(db *sqlx.DB, removeReportInfo *RemoveReportInfo) error {

	tranX, err := db.Beginx()
	if err != nil {
		return err
	}

	defer tranX.Rollback()

	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1 FROM "Reports"
			WHERE report_id = $1
		)
	`
	err = tranX.Get(&exists, query, removeReportInfo.Report_id)
	if err != nil {
		return err
	}

	if exists {
		_, err = tranX.Exec(
			`DELETE FROM "Reports"
			WHERE report_id = $1
			`, removeReportInfo.Report_id)
		if err != nil {
			return err
		}
	} else {
		return errors.New("ไม่พบรายงานนี้")
	}

	err = tranX.Commit()
	if err != nil {
		return err
	}

	return nil
}
