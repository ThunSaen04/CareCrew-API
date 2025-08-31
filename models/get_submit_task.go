package models

import "github.com/jmoiron/sqlx"

type SubmitTaskWithID struct {
	PersonnelID int    `db:"personnel_id" json:"personnel_id"`
	Task_id     int    `db:"task_id" json:"task_id"`
	Submit      bool   `db:"submit" json:"submit"`
	Submit_at   string `db:"submit_at" json:"submit_at"`
}

func Get_Submit_Task_With_ID(db *sqlx.DB, Personnel_id int, Task_id int) (*SubmitTaskWithID, error) {
	var smtw_id SubmitTaskWithID
	query := `
        SELECT tas.personnel_id, tas.task_id, tas.submit, tas.submit_at
        FROM "Tasks_assignment" tas
		WHERE tas.personnel_id = $1 AND tas.task_id = $2
    `
	err := db.Get(&smtw_id, query, Personnel_id, Task_id)
	if err != nil {
		return nil, err
	}
	return &smtw_id, nil
}
