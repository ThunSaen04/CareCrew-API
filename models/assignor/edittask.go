package assignor

import (
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/project/carecrew/orther"
)

type EditTaskInfo struct {
	Personnel_id     int    `db:"personnel_id" json:"personnel_id"`
	Task_id          int    `db:"task_id" json:"task_id"`                   //Tasks
	Task_type_id     int    `db:"task_type_id" json:"task_type_id"`         //Tasks
	Title            string `db:"title" json:"title"`                       //Tasks_detail
	Detail           string `db:"detail" json:"detail"`                     //Tasks_detail
	Location         string `db:"location" json:"location"`                 //Tasks_detail
	Priority_type_id int    `db:"priority_type_id" json:"priority_type_id"` //Tasks
	People_needed    int    `db:"people_needed" json:"people_needed"`       //Tasks_detail
	Task_due_at      string `db:"task_due_at" json:"task_due_at"`           //Task
}

func EditTask(db *sqlx.DB, edittaskinfo *EditTaskInfo) error {
	tranX, err := db.Beginx()
	if err != nil {
		return err
	}

	defer tranX.Rollback()

	//เช็คว่ามีงานนั้นจริงอ๊ะป่าว
	var d []string
	var exists bool
	query := `
			SELECT EXISTS (
				SELECT 1 FROM "Tasks"
				WHERE task_id = $1
			)
		`
	err = tranX.Get(&exists, query, edittaskinfo.Task_id)
	if err != nil {
		return err
	}

	if exists { //ถ้ามีงานจริง
		//ข้อมูลงาน
		var taskinfo EditTaskInfo
		query = `
		SELECT t.task_id, t.task_type_id, td.title, td.detail, td.location, t.priority_type_id, td.people_needed, t.task_due_at
		FROM "Tasks" t
		LEFT JOIN "Tasks_detail" td	ON t.task_id = td.task_id
		WHERE t.task_id = $1
	`
		err = tranX.Get(&taskinfo, query, edittaskinfo.Task_id)
		if err != nil {
			return errors.New("เกิดข้อผิดพลาดในการเรียกข้อมูลงาน")
		}

		//กำหนดส่งงาน
		if edittaskinfo.Task_due_at != "" && edittaskinfo.Task_due_at != taskinfo.Task_due_at {
			query = `
				UPDATE "Tasks"
				SET task_due_at = $1
				WHERE task_id = $2
			`
			_, err = tranX.Exec(query, edittaskinfo.Task_due_at, edittaskinfo.Task_id)
			if err != nil {
				return errors.New("แก้ไข กำหนดส่งงาน ไม่สำเร็จ")
			}
			d = append(d, "แก้ไขกำหนดส่งงาน")
		}

		//ประเภทงาน
		if edittaskinfo.Task_type_id != 0 && edittaskinfo.Task_type_id != taskinfo.Task_type_id {
			query = `
				UPDATE "Tasks"
				SET task_type_id = $1
				WHERE task_id = $2
			`
			_, err = tranX.Exec(query, edittaskinfo.Task_type_id, edittaskinfo.Task_id)
			if err != nil {
				return errors.New("แก้ไข ประเภทงาน ไม่สำเร็จ")
			}
			d = append(d, "แก้ไขประเภทงาน")
		}
		//หัวข้องาน
		if len(edittaskinfo.Title) != 0 && edittaskinfo.Title != taskinfo.Title {
			query = `
				UPDATE "Tasks_detail"
				SET title = $1
				WHERE task_id = $2
				
			`
			_, err = tranX.Exec(query, edittaskinfo.Title, edittaskinfo.Task_id)
			if err != nil {
				return errors.New("แก้ไข หัวข้องาน ไม่สำเร็จ")
			}
			d = append(d, "แก้ไขหัวข้องาน")
		}
		//รายละเอียดงาน
		if len(edittaskinfo.Detail) != 0 && edittaskinfo.Detail != taskinfo.Detail {
			query = `
				UPDATE "Tasks_detail"
				SET detail = $1
				WHERE task_id = $2
			`
			_, err = tranX.Exec(query, edittaskinfo.Detail, edittaskinfo.Task_id)
			if err != nil {
				return errors.New("แก้ไข รายละเอียดงาน ไม่สำเร็จ")
			}
			d = append(d, "แก้ไขรายละเอียดงาน")
		}
		//สถานที่งาน
		if len(edittaskinfo.Location) != 0 && edittaskinfo.Location != taskinfo.Location {
			query = `
				UPDATE "Tasks_detail"
				SET location = $1
				WHERE task_id = $2
			`
			_, err = tranX.Exec(query, edittaskinfo.Location, edittaskinfo.Task_id)
			if err != nil {
				return errors.New("แก้ไข สถานที่งาน ไม่สำเร็จ")
			}
			d = append(d, "แก้ไขสถานที่งาน")
		}
		//ความสำคัญงงาน
		if edittaskinfo.Priority_type_id != 0 && edittaskinfo.Priority_type_id != taskinfo.Priority_type_id {
			query = `
				UPDATE "Tasks"
				SET priority_type_id = $1
				WHERE task_id = $2
			`
			_, err = tranX.Exec(query, edittaskinfo.Priority_type_id, edittaskinfo.Task_id)
			if err != nil {
				return errors.New("แก้ไข ความสำคัญงาน ไม่สำเร็จ")
			}
			d = append(d, "แก้ไขความสำคัญงาน")
		}
		//จำนวนบุคลากรที่ต้องการ
		if edittaskinfo.People_needed != 0 && edittaskinfo.People_needed != taskinfo.People_needed {
			query = `
				UPDATE "Tasks_detail"
				SET people_needed = $1
				WHERE task_id = $2
			`
			_, err = tranX.Exec(query, edittaskinfo.People_needed, edittaskinfo.Task_id)
			if err != nil {
				return errors.New("แก้ไข จำนวนบุคลากรที่ต้องการ ไม่สำเร็จ")
			}
			d = append(d, "แก้ไขจำนวนบุคลากรที่ต้องการ")
		}

		query = `
			UPDATE "Tasks"
			SET updated_at = NOW()
			WHERE task_id = $1
		`
		_, err = tranX.Exec(query, edittaskinfo.Task_id)
		if err != nil {
			return err
		}

	} else {
		return errors.New("ไม่พบงานดังกล่าว")
	}

	_, err = tranX.Exec(`
			INSERT INTO "Assignor_logs" (personnel_id, task_id, detail)
			VALUES ($1, $2, $3)
		`, edittaskinfo.Personnel_id, edittaskinfo.Task_id, "แก้ไขงาน")
	if err != nil {
		return err
	}

	err = tranX.Commit()
	if err != nil {
		return err
	}

	sendinfo := orther.SendNotiInfo{
		Task_id: edittaskinfo.Task_id,
		Detail:  strings.Join(d, ", "),
		Title:   "งาน " + edittaskinfo.Title,
		Body:    "มีการแก้ไขรายละเอียดงาน",
	}

	orther.SendNotiSuccessToPerInTask(db, &sendinfo) //

	return nil
}
