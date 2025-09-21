package worker

import (
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
)

type PerLrubTaskInfo struct {
	Personnel_id int `db:"personnel_id" json:"personnel_id"` //Personnels
	Task_id      int `db:"task_id" json:"task_id"`           //Tasks
}

type taskinfoget struct {
	Status_Type_id  int `db:"status_type_id" json:"status_type_id"`
	People_Needed   int `db:"people_needed" json:"people_needed"`
	Personnel_Count int `db:"personnel_count" json:"personnel_count"`
}

func PerLrubTask(db *sqlx.DB, perLrubTaskInfo *PerLrubTaskInfo) error {
	tranX, err := db.Beginx()
	if err != nil {
		return err
	}

	defer tranX.Rollback()

	/////////////////////////////////////////////////////////////////////////////////////////////////////
	var exists bool
	//เช็ครับงานซ้ำ
	query := `
		SELECT EXISTS (
			SELECT 1 FROM "Tasks_assignment"
			WHERE task_id = $1 AND personnel_id = $2
		)
	`
	err = tranX.Get(&exists, query, perLrubTaskInfo.Task_id, perLrubTaskInfo.Personnel_id)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("ผู้ใช้งานหมายเลขนี้ได้รับงานนี้แล้ว")
	}
	//log.Print("[System] ผ่าน 1 เช็ครับงานซ้ำ")

	/////////////////////////////////////////////////////////////////////////////////////////////////////
	var task_id int
	// ผู้ใช้งาน รับงาน
	err = tranX.QueryRow(
		`INSERT INTO "Tasks_assignment" (personnel_id, task_id, accecp_at)
		VALUES ($1, $2, NOW())
		RETURNING task_id`,
		perLrubTaskInfo.Personnel_id,
		perLrubTaskInfo.Task_id,
	).Scan(&task_id)
	if err != nil {
		return err
	}
	//log.Print("[System] ผ่าน 2 ผู้ใช้งาน รับงาน")

	/////////////////////////////////////////////////////////////////////////////////////////////////////
	var taskinfoget taskinfoget
	// ตรวจสอบ สถานะงาน และเปลี่ยนสถานะงาน / ตรวจสอบจำนวนคนที่รับงาน
	query = `
		SELECT t.status_type_id, td.people_needed, COUNT(ta.personnel_id) AS personnel_count
		FROM "Tasks" t
		LEFT JOIN "Tasks_detail" td ON t.task_id = td.task_id
		LEFT JOIN "Tasks_assignment" ta ON t.task_id = ta.task_id
		WHERE t.task_id = $1
		GROUP BY t.status_type_id, td.people_needed, t.task_id
		`
	err = tranX.Get(&taskinfoget, query, task_id)
	if err != nil {
		return err
	}
	//log.Print("[System] ผ่าน 3 เรียกข้อมูล")

	if taskinfoget.Personnel_Count-1 < taskinfoget.People_Needed {
		_, err := tranX.Exec(
			`UPDATE "Tasks" 
				SET updated_at = NOW()
				WHERE task_id = $1
				`,
			task_id,
		)
		if err != nil {
			return err
		}
		//log.Print("[System] ผ่าน 3.2 updated_at = NOW()")

		if taskinfoget.Status_Type_id == 3 {
			_, err := tranX.Exec(
				`UPDATE "Tasks" 
				SET status_type_id = 2, updated_at = NOW()
				WHERE task_id = $1
				`,
				task_id,
			)
			if err != nil {
				return err
			}
			//log.Print("[System] ผ่าน 3.3 status_type_id = 2, updated_at = NOW()")

		} else if taskinfoget.Status_Type_id == 1 {
			return errors.New("งานนี้เสร็จสิ้นแล้ว")
		} else if taskinfoget.Status_Type_id == 4 {
			return errors.New("งานนี้อยู่ระหว่างการตรวจสอบ")
		}
	} else {
		return errors.New("จำนวนคนรับงานครบแล้ว")
	}

	_, err = tranX.Exec(`
			INSERT INTO "Worker_logs" (personnel_id, task_id, detail)
			VALUES ($1, $2, $3)
		`, perLrubTaskInfo.Personnel_id, perLrubTaskInfo.Task_id, "รับงาน")
	if err != nil {
		return err
	}

	/////////////////////////////////////////////////////////////////////////////////////////////////////
	err = tranX.Commit()
	if err != nil {
		return err
	}

	log.Print("[System] ผู้ใช้งานหมายเลข:", perLrubTaskInfo.Personnel_id, "ได้รับงาน: ", perLrubTaskInfo.Task_id)
	return nil
}
