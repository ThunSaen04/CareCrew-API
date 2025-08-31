package models

import "github.com/jmoiron/sqlx"

type TasktypeInfo struct {
	Task_type_id   int    `db:"task_type_id" json:"task_type_id"` //Tasks_detail
	Task_type_name string `db:"name" json:"name"`                 //Tasks_detail
}

func Get_Task_Type_Info(db *sqlx.DB) ([]TasktypeInfo, error) {
	tasktypeinfo := []TasktypeInfo{}
	query := `
		SELECT tt.task_type_id, tt.name
		FROM "Task_types" tt
	`
	err := db.Select(&tasktypeinfo, query)
	return tasktypeinfo, err
}
