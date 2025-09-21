package assignor

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type RemoveTaskInfo struct {
	Personnel_id int `db:"personnel_id" json:"personnel_id"`
	Task_id      int `db:"task_id" json:"task_id"`
}

func RemoveTask(db *sqlx.DB, removeTaskInfo *RemoveTaskInfo) error {

	tranX, err := db.Beginx()
	if err != nil {
		return err
	}

	defer tranX.Rollback()

	var exists bool
	//เช็คว่ามีงานนั้นจริงอ๊ะป่าว
	query := `
		SELECT EXISTS (
			SELECT 1 FROM "Tasks"
			WHERE task_id = $1
		)
	`
	err = tranX.Get(&exists, query, removeTaskInfo.Task_id)
	if err != nil {
		return err
	}

	if exists {
		_, err = tranX.Exec(
			`DELETE FROM "Tasks"
			WHERE task_id = $1
			`, removeTaskInfo.Task_id)
		if err != nil {
			return err
		}
	} else {
		return errors.New("ไม่พบงาน")
	}

	_, err = tranX.Exec(`
			INSERT INTO "Assignor_logs" (personnel_id, task_id, detail)
			VALUES ($1, $2, $3)
		`, removeTaskInfo.Personnel_id, removeTaskInfo.Task_id, "ลบงาน")
	if err != nil {
		return err
	}

	err = tranX.Commit()
	if err != nil {
		return err
	}

	return nil
}
