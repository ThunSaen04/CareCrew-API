package worker

import (
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
)

type YokLerkTaskInfo struct {
	Personnel_id int `db:"personnel_id" json:"personnel_id"`
	Task_id      int `db:"task_id" json:"task_id"`
}

func YokLerkTask(db *sqlx.DB, yokLerkTaskInfo *YokLerkTaskInfo) error {

	tranX, err := db.Beginx()
	if err != nil {
		return err
	}

	defer tranX.Rollback()

	var exists bool
	//เช็คว่ามีงานนั้นจริงอ๊ะป่าว
	query := `
		SELECT EXISTS (
			SELECT 1 FROM "Tasks_assignment"
			WHERE task_id = $1 AND personnel_id = $2
		)
	`
	err = tranX.Get(&exists, query, yokLerkTaskInfo.Task_id, yokLerkTaskInfo.Personnel_id)
	if err != nil {
		return err
	}
	if exists {
		_, err = tranX.Exec(`
		DELETE FROM "Tasks_assignment"
		WHERE task_id = $1 AND personnel_id = $2
		`, yokLerkTaskInfo.Task_id, yokLerkTaskInfo.Personnel_id)
		if err != nil {
			return errors.New("ไม่สามารถยกเลิกงานได้")
		}
	} else {
		return errors.New("ผู้ใช้งานไม่ได้รับงานดังกล่าว")
	}

	var taskinfoget taskinfoget
	query = `
		SELECT t.status_type_id, COUNT(ta.personnel_id) AS personnel_count
		FROM "Tasks" t
		LEFT JOIN "Tasks_assignment" ta ON t.task_id = ta.task_id
		WHERE t.task_id = $1
		GROUP BY t.status_type_id, t.task_id
		`
	err = tranX.Get(&taskinfoget, query, yokLerkTaskInfo.Task_id)
	if err != nil {
		return err
	}

	if taskinfoget.Personnel_Count == 0 && taskinfoget.Status_Type_id == 2 {
		_, err := tranX.Exec(
			`UPDATE "Tasks" 
				SET status_type_id = 3, updated_at = NOW()
				WHERE task_id = $1
				`,
			yokLerkTaskInfo.Task_id,
		)
		if err != nil {
			return err
		}
	} else if taskinfoget.Status_Type_id == 1 {
		return errors.New("ไม่สามารถยกเลิกงานได้")
	} else if taskinfoget.Status_Type_id == 4 {
		return errors.New("ไม่สามารถยกเลิกงานได้ กรุณายกเลิกการส่งงานก่อน")
	}

	err = tranX.Commit()
	if err != nil {
		return err
	}

	log.Print("[System] ผู้ใช้งานหมายเลข:", yokLerkTaskInfo.Personnel_id, "ได้ยกเลิกงาน: ", yokLerkTaskInfo.Task_id)
	return nil
}
