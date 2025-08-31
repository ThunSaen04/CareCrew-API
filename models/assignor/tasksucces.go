package assignor

import (
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
)

type taskStatus struct {
	Completed bool `db:"completed" json:"completed"`
	Status    int  `db:"status_type_id" json:"status_type_id"`
}

func TaskSuccess(db *sqlx.DB, taskID int) error {
	tranX, err := db.Beginx()
	if err != nil {
		return err
	}

	defer tranX.Rollback()

	var task taskStatus
	query := `
		SELECT t.completed, t.status_type_id
		FROM "Tasks" t
		WHERE t.task_id = $1
	`
	err = tranX.Get(&task, query, taskID)
	if err != nil {
		log.Printf("[Warning] ไม่พบงานหมายเลข %d", taskID)
		return err
	}

	if !task.Completed && task.Status == 4 {
		_, err = tranX.Exec(
			`UPDATE "Tasks"
			SET status_type_id = 1, updated_at = NOW(), completed = true, completed_at = NOW()
			WHERE task_id = $1`,
			taskID,
		)
		if err != nil {
			return err
		}
	} else {
		log.Print("[Warning] สถานะงานไม่ถูกต้อง หรืองานที่สิ้นสุดแล้ว")
		return errors.New("สถานะงานไม่ถูกต้อง หรืองานที่สิ้นสุดแล้ว")
	}

	err = tranX.Commit()
	if err != nil {
		return err
	}

	log.Printf("[System] ยืนยันการตรวบสอบงานหมายเลข: %d แล้ว", taskID)
	return nil
}
