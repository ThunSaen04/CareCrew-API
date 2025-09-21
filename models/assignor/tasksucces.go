package assignor

import (
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/project/carecrew/orther"
)

type taskStatus struct {
	Personnel_id int    `db:"personnel_id" json:"personnel_id"`
	Completed    bool   `db:"completed" json:"completed"`
	Status       int    `db:"status_type_id" json:"status_type_id"`
	Title        string `db:"title" json:"title"`
}

func TaskSuccess(db *sqlx.DB, taskID int) error {
	tranX, err := db.Beginx()
	if err != nil {
		return err
	}

	defer tranX.Rollback()

	var task taskStatus
	query := `
		SELECT t.completed, t.status_type_id, td.title
		FROM "Tasks" t
		LEFT JOIN "Tasks_detail" td ON t.task_id = td.task_id
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

	_, err = tranX.Exec(`
			INSERT INTO "Assignor_logs" (personnel_id, task_id, detail)
			VALUES ($1, $2, $3)
		`, task.Personnel_id, taskID, "อนุมัติเสร็จสิ้นงาน")
	if err != nil {
		return err
	}

	err = tranX.Commit()
	if err != nil {
		return err
	}
	sendinfo := orther.SendNotiInfo{
		Task_id: taskID,
		Title:   "งาน " + task.Title + " สิ้นสุดแล้ว!!",
		Body:    "งานนี้สิ้นสุดเรียบร้อย",
	}

	orther.SendNotiSuccessToPerInTask(db, &sendinfo)
	log.Printf("[System] ยืนยันการตรวบสอบงานหมายเลข: %d แล้ว", taskID)
	return nil
}
