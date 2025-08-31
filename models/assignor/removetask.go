package assignor

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type RemoveTaskInfo struct {
	Task_id int `db:"task_id" json:"task_id"`
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

	err = tranX.Commit()
	if err != nil {
		return err
	}

	return nil
}
