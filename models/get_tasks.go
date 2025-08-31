package models

import (
	"github.com/jmoiron/sqlx"
)

type TasksInfo struct {
	Task_id       int    `db:"task_id" json:"task_id"`             //Tasks
	Task_type     string `db:"type_name" json:"type_name"`         //Task_types
	Task_title    string `db:"title" json:"title"`                 //Tasks_detail
	Detail        string `db:"detail" json:"detail"`               //Tasks_detail
	Location      string `db:"location" json:"location"`           //Tasks_detail
	People_needed int    `db:"people_needed" json:"people_needed"` //Tasks_detail
	Assigned_by   string `db:"assigned_by" json:"assigned_by"`     //Personnels
	Priority_type string `db:"priority_name" json:"priority_name"` //Priority_type
	Status_types  string `db:"status_type" json:"status"`          //Status_types
	Created_at    string `db:"created_at" json:"created_at"`       //Tasks
	Updated_at    string `db:"updated_at" json:"updated_at"`       //Tasks
}

func GetTasks(db *sqlx.DB) ([]TasksInfo, error) {
	tasks := []TasksInfo{}
	query := `
        SELECT t.task_id, tt.name AS type_name, td.title, td.detail, td.location, pt.name AS priority_name, td.people_needed, p.first_name || ' ' || p.last_name AS assigned_by, st.name AS status_type, t.created_at, t.updated_at
        FROM "Tasks" t
        LEFT JOIN "Task_types" tt ON t.task_type_id = tt.task_type_id
        LEFT JOIN "Tasks_detail" td ON t.task_id = td.task_id
        LEFT JOIN "Personnels" p ON td.assigned_by = p.personnel_id
        LEFT JOIN "Status_types" st ON t.status_type_id = st.status_type_id
        LEFT JOIN "Priority_types" pt ON t.priority_type_id = pt.priority_type_id
        ORDER BY st.status_type_id DESC, t.priority_type_id, t.task_id DESC
    `
	err := db.Select(&tasks, query)
	return tasks, err
}
