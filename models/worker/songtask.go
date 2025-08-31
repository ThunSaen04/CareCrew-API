package worker

import (
	"errors"

	"github.com/jmoiron/sqlx"
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

	defer tranX.Rollback()

	//เช็คว่ามีงานนั้นจริงอ๊ะป่าว หรือว่าส่งงานไปยัง???
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1 FROM "Tasks_assignment"
			WHERE task_id = $1 AND personnel_id = $2 AND submit = false
		)
	`
	err = tranX.Get(&exists, query, persongTaskInfo.Task_id, persongTaskInfo.Personnel_id)
	if err != nil {
		return err
	}
	if exists {
		//อัพเดท submit ให้เป็น True
		var assignment_id int
		err = tranX.QueryRow(`
			UPDATE "Tasks_assignment"
			SET submit = true, submit_at = NOW()
			WHERE personnel_id = $1 AND task_id = $2
			RETURNING assignment_id
		`, persongTaskInfo.Personnel_id, persongTaskInfo.Task_id).Scan(&assignment_id)
		if err != nil {
			return err
		}

		//ใส่ข้อมูลในการส่งงาน
		_, err = tranX.Exec(`
			INSERT INTO "Tasks_attachments" (assignment_id, file, uploaded_at)
			VALUES ($1, $2, NOW())
		`, assignment_id, persongTaskInfo.File)
		if err != nil {
			return err
		}

		//ตรวจสอบว่าส่งครบยัง ถ้าส่งครบทุกคนก็เปลี่ยนสถานะงานซะ
		var info submitTaskInfo
		query := `
			SELECT 
				COUNT(personnel_id) AS personnel_count,
				SUM(CASE WHEN submit = true THEN 1 ELSE 0 END) AS submitted_count
			FROM "Tasks_assignment"
			WHERE task_id = $1
			`
		err := tranX.Get(&info, query, persongTaskInfo.Task_id)
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
		}

	} else {
		return errors.New("ไม่พบงานที่รับหรือพบว่ามีการส่งงานนี้แล้ว")
	}

	err = tranX.Commit()
	if err != nil {
		return err
	}

	return nil
}
