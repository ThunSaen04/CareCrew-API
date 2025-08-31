package v2

import (
	"github.com/jmoiron/sqlx"
	"github.com/project/carecrew/models"
)

func LrubTasksCountV2(db *sqlx.DB, task_id int) ([]models.LrubTasksCountInfo, error) {
	var lrubAs []models.LrubTasksCountInfo
	query := `
		SELECT ta.task_id, STRING_AGG(p.personnel_id::TEXT, ', ') AS personnel_ids, COUNT(ta.personnel_id) AS personnel_count, STRING_AGG(p.first_name || ' ' || p.last_name, ', ') AS personnel_name
		FROM "Tasks_assignment" ta
		INNER JOIN "Personnels" p ON ta.personnel_id = p.personnel_id
		WHERE task_id = $1
		GROUP BY ta.task_id;
	`
	err := db.Select(&lrubAs, query, task_id)
	if err != nil {
		return nil, err
	}
	return lrubAs, nil
}
