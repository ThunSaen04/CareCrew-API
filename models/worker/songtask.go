package worker

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/project/carecrew/orther"
)

type PerSongTaskInfo struct {
	Personnel_id int    `db:"personnel_id" json:"personnel_id"`
	Task_id      int    `db:"task_id" json:"task_id"`
	File         string `db:"file" json:"file"`
}

type submitTaskInfo struct {
	PersonnelCount int `db:"personnel_count"`
	SubmittedCount int `db:"submitted_count"`
}

func Songtask(db *sqlx.DB, persongTaskInfo *PerSongTaskInfo) error {

	tranX, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tranX.Rollback()

	var assignment_id int
	query := `
		SELECT assignment_id, submit
		FROM "Tasks_assignment"
		WHERE task_id = $1 AND personnel_id = $2
	`
	var submit bool
	err = tranX.QueryRow(query, persongTaskInfo.Task_id, persongTaskInfo.Personnel_id).
		Scan(&assignment_id, &submit)
	if err != nil {
		return errors.New("ไม่พบงานที่รับไว้")
	}

	if !submit {
		_, err = tranX.Exec(`
			UPDATE "Tasks_assignment"
			SET submit = true, submit_at = NOW()
			WHERE assignment_id = $1
		`, assignment_id)
		if err != nil {
			return err
		}

		_, err = tranX.Exec(`
			INSERT INTO "Tasks_attachments" (assignment_id, file, uploaded_at)
			VALUES ($1, $2, NOW())
		`, assignment_id, persongTaskInfo.File)
		if err != nil {
			return err
		}

		var info submitTaskInfo
		query := `
			SELECT 
				COUNT(personnel_id) AS personnel_count,
				SUM(CASE WHEN submit = true THEN 1 ELSE 0 END) AS submitted_count
			FROM "Tasks_assignment"
			WHERE task_id = $1
		`
		err = tranX.Get(&info, query, persongTaskInfo.Task_id)
		if err != nil {
			return err
		}

		if info.PersonnelCount == info.SubmittedCount {
			_, err = tranX.Exec(`
				UPDATE "Tasks"
				SET status_type_id = 4, updated_at = NOW()
				WHERE task_id = $1`,
				persongTaskInfo.Task_id)
			if err != nil {
				return err
			}
			sendinfo := orther.SendNotiInfo{
				Task_id: persongTaskInfo.Task_id,
				Detail:  time.Now().Format("2006-01-02 15:04:05-07"),
				Title:   "งานรอการตรวจสอบ!!",
				Body:    "พบงานที่ต้องการ การตรวจสอบใหม่",
			}
			orther.SendNotiSuccessToAssignor(db, &sendinfo) //
		}

	} else {
		if persongTaskInfo.File != "" {
			_, err = tranX.Exec(`
				UPDATE "Tasks_attachments"
				SET file = $1, uploaded_at = NOW()
				WHERE assignment_id = $2
			`, persongTaskInfo.File, assignment_id)
			if err != nil {
				return err
			}
		}
	}

	_, err = tranX.Exec(`
		INSERT INTO "Worker_logs" (personnel_id, task_id, detail, file)
		VALUES ($1, $2, $3, $4)
	`, persongTaskInfo.Personnel_id, persongTaskInfo.Task_id, "ส่งงาน", persongTaskInfo.File)
	if err != nil {
		return err
	}

	err = tranX.Commit()
	if err != nil {
		return err
	}

	return nil
}
