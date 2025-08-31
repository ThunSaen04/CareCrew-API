package models

import (
	"github.com/jmoiron/sqlx"
)

type PersonnelLrubTaskInfo struct {
	PersonnelID int    `db:"personnel_id" json:"personnel_id"`
	Task_id     string `db:"task_id" json:"task_id"` // task_id รวมกันเป็นสตริง เช่น "1,3,5"
}

func PersonnelLrubTask(db *sqlx.DB) ([]PersonnelLrubTaskInfo, error) {
	var lrubAs []PersonnelLrubTaskInfo
	query := `
        SELECT personnel_id, STRING_AGG(task_id::text, ',') AS task_id
        FROM "Tasks_assignment"
        GROUP BY personnel_id;
	`

	err := db.Select(&lrubAs, query)
	if err != nil {
		return nil, err
	}
	return lrubAs, nil
}
