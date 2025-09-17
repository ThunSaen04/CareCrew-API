package models

import "github.com/jmoiron/sqlx"

type TaskpriorityInfo struct {
	Task_Priority_id   int    `db:"priority_type_id" json:"priority_type_id"` //Priority_types
	Task_Priority_name string `db:"name" json:"name"`                         //Priority_types
}

func Get_Task_Priority_Info(db *sqlx.DB) ([]TaskpriorityInfo, error) {
	taskpriorityinfo := []TaskpriorityInfo{}
	query := `
		SELECT pt.priority_type_id, pt.name
		FROM "Priority_types" pt
		ORDER BY pt.priority_type_id
	`
	err := db.Select(&taskpriorityinfo, query)
	return taskpriorityinfo, err
}
