package assignor

import (
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/project/carecrew/orther"
)

type NoSuccessInfo struct {
	Task_id      int    `db:"task_id" json:"task_id"`
	Personnel_id int    `db:"personnel_id" json:"personnel_id"`
	Detail       string `db:"detail" json:"detail"`
}

func NoSuccess(db *sqlx.DB, nosuccessinfo NoSuccessInfo) error {
	tranX, err := db.Beginx()
	if err != nil {
		return err
	}

	defer tranX.Rollback()

	var exists bool
	query := `
			SELECT EXISTS (
				SELECT 1 FROM "Tasks"
				WHERE task_id = $1
			)
		`
	err = tranX.Get(&exists, query, nosuccessinfo.Task_id)
	if err != nil {
		return err
	}

	if exists {
		var task taskStatus
		query = `
			SELECT t.completed, t.status_type_id, td.title
			FROM "Tasks" t
			LEFT JOIN "Tasks_detail" td ON t.task_id = td.task_id
			WHERE t.task_id = $1
		`
		err = tranX.Get(&task, query, nosuccessinfo.Task_id)
		if err != nil {
			log.Printf("[Warning] ไม่พบงานหมายเลข %d", nosuccessinfo.Task_id)
			return err
		}

		if task.Completed == false && task.Status == 4 {
			query = `
				UPDATE "Tasks"
				SET status_type_id = 5, completed = false, completed_at = NULL
				WHERE task_id = $1
			`
			_, err = tranX.Exec(query, nosuccessinfo.Task_id)
			if err != nil {
				return err
			}

			query = `
				UPDATE "Tasks_assignment"
				SET submit = false, submit_at = NULL
				WHERE task_id = $1
				RETURNING assignment_id
			`
			rows, err := tranX.Query(query, nosuccessinfo.Task_id)
			if err != nil {
				return err
			}
			defer rows.Close()

			var assignmentIDs []int
			for rows.Next() {
				var id int
				err := rows.Scan(&id)
				if err != nil {
					return err
				}
				assignmentIDs = append(assignmentIDs, id)
			}

			for _, id := range assignmentIDs {
				dquery := `DELETE FROM "Tasks_attachments" WHERE assignment_id = $1`
				_, err = tranX.Exec(dquery, id)
				if err != nil {
					return err
				}
			}

			query = `
				UPDATE "Tasks"
				SET nosuccess_detail = $2, updated_at = now()
				WHERE task_id = $1
			`
			_, err = tranX.Exec(query, nosuccessinfo.Task_id, nosuccessinfo.Detail)
			if err != nil {
				return err
			}
		} else {
			return errors.New("สถานะงานไม่อยู่ในเงื่อนไข")
		}

		_, err = tranX.Exec(`
			INSERT INTO "Assignor_logs" (personnel_id, task_id, detail)
			VALUES ($1, $2, $3)
		`, nosuccessinfo.Personnel_id, nosuccessinfo.Task_id, "ไม่อนุมัติต้องแก้ไขงาน")
		if err != nil {
			return err
		}

		err = tranX.Commit()
		if err != nil {
			return err
		}

		sendinfo := orther.SendNotiInfo{
			Task_id: nosuccessinfo.Task_id,
			Detail:  nosuccessinfo.Detail,
			Title:   "งาน " + task.Title,
			Body:    "ไม่อนุมัติ กรุณาแก้ไขงาน",
		}
		orther.SendNotiSuccessToPerInTask(db, &sendinfo) //
	} else {
		return errors.New("ไม่พบงาน")
	}

	return nil
}
