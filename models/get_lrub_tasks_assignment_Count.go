package models

import (
	"github.com/jmoiron/sqlx"
)

type LrubTasksCountInfo struct {
	Task_id         int    `db:"task_id" json:"task_id"`                 //Tasks_assignment
	Personanel_Name string `db:"personnel_name" json:"personnel_name"`   //Personnels
	Personanel_ID   string `db:"personnel_ids" json:"personnel_ids"`     //Personnels
	PersonnelCount  int    `db:"personnel_count" json:"personnel_count"` //Personnels
}

func LrubTasksCount(db *sqlx.DB) ([]LrubTasksCountInfo, error) {
	var lrubAs []LrubTasksCountInfo
	query := `
		SELECT ta.task_id, STRING_AGG(p.personnel_id::TEXT, ', ') AS personnel_ids, COUNT(ta.personnel_id) AS personnel_count, STRING_AGG(p.first_name || ' ' || p.last_name, ', ') AS personnel_name
		FROM "Tasks_assignment" ta
		INNER JOIN "Personnels" p ON ta.personnel_id = p.personnel_id
		GROUP BY ta.task_id;
	`

	err := db.Select(&lrubAs, query)
	if err != nil {
		return nil, err
	}
	return lrubAs, nil
}
