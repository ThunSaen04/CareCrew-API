package assignor

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type RemoveReportInfo struct {
	Personnel_id int `db:"personnel_id" json:"personnel_id"`
	Report_id    int `db:"report_id" json:"report_id"`
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

	_, err = tranX.Exec(`
			INSERT INTO "Assignor_logs" (personnel_id, report_id, detail)
			VALUES ($1, $2, $3)
		`, removeReportInfo.Personnel_id, removeReportInfo.Report_id, "ลบการแจ้งเหตุเพื่อสร้างงาน")
	if err != nil {
		return err
	}

	err = tranX.Commit()
	if err != nil {
		return err
	}

	return nil
}
