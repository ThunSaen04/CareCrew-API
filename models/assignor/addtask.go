package assignor

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/project/carecrew/orther"
)

type AddTaskInfo struct {
	Task_type_id     int    `db:"task_type_id" json:"task_type_id"`         //Task
	Priority_type_id int    `db:"priority_type_id" json:"priority_type_id"` //Task
	Title            string `db:"title" json:"title"`                       //Tasks_detail
	Detail           string `db:"detail" json:"detail"`                     //Tasks_detail
	Location         string `db:"location" json:"location"`                 //Tasks_detail
	People_needed    int    `db:"people_needed" json:"people_needed"`       //Tasks_detail
	Assigned_by      int    `db:"assigned_by" json:"assigned_by"`           //Tasks_detail (ใช้ PersonnelIDDD)
}

func AddTask(db *sqlx.DB, addnewtaskinfo *AddTaskInfo) error {
	tranX, err := db.Beginx()
	if err != nil {
		return err
	}

	defer tranX.Rollback()

	if addnewtaskinfo.People_needed <= 0 {
		return errors.New("กรุณาระบุจำนวนบุคลากรที่ต้องการ (ขั้นต่ำ 1คน)")
	}

	var task_id int

	err = tranX.QueryRow(
		`INSERT INTO "Tasks" (task_type_id, priority_type_id)
		VALUES ($1, $2)
		RETURNING task_id`,
		addnewtaskinfo.Task_type_id,
		addnewtaskinfo.Priority_type_id,
	).Scan(&task_id)
	if err != nil {
		return err
	}

	_, err = tranX.Exec(
		`INSERT INTO "Tasks_detail" (task_id, title, detail, location, people_needed, assigned_by)
			VALUES ($1, $2, $3, $4, $5, $6)`,
		task_id,
		addnewtaskinfo.Title,
		addnewtaskinfo.Detail,
		addnewtaskinfo.Location,
		addnewtaskinfo.People_needed,
		addnewtaskinfo.Assigned_by,
	)
	if err != nil {
		return err
	}

	err = tranX.Commit()
	if err != nil {
		return err
	}
	orther.SendNotificationToAll(db, "งานใหม่มาแล้ว!!!", addnewtaskinfo.Title)

	return err
}
