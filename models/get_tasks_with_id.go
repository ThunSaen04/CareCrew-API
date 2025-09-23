package models

import (
	"github.com/jmoiron/sqlx"
)

func GetTasks_With_ID(db *sqlx.DB, TasksID int) (*TasksInfo, error) {
	var tasks TasksInfo
	query := `
        SELECT t.task_id, tt.name AS type_name, td.title, td.detail, td.location, pt.name AS priority_name, td.people_needed, p.first_name || ' ' || p.last_name AS assigned_by, st.name AS status_type, t.created_at, t.updated_at, t.task_due_at, COUNT(ta.personnel_id) AS personnel_count, t.nosuccess_detail
        FROM "Tasks" t
        LEFT JOIN "Task_types" tt ON t.task_type_id = tt.task_type_id
        LEFT JOIN "Tasks_detail" td ON t.task_id = td.task_id
        LEFT JOIN "Personnels" p ON td.assigned_by = p.personnel_id
        LEFT JOIN "Status_types" st ON t.status_type_id = st.status_type_id
        LEFT JOIN "Priority_types" pt ON t.priority_type_id = pt.priority_type_id
		LEFT JOIN "Tasks_assignment" ta ON ta.task_id = t.task_id
		WHERE t.task_id = $1
		GROUP BY 
			t.task_id, tt.name, td.title, td.detail, td.location, pt.name, 
			td.people_needed, p.first_name, p.last_name, st.name, st.status_type_id,
			t.created_at, t.updated_at, t.task_due_at, t.nosuccess_detail
		ORDER BY st.status_type_id DESC, t.priority_type_id, t.task_id DESC;
    `

	err := db.Get(&tasks, query, TasksID)
	if err != nil {
		return nil, err
	}
	return &tasks, nil
}
