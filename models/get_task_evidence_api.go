package models

import (
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/project/carecrew/config"
)

type TaskEvidenceInfo struct {
	Task_id       int      `db:"task_id" json:"task_id"`             //Tasks_assignment
	Task_title    string   `db:"title" json:"title"`                 //Tasks_detail
	Detail        string   `db:"detail" json:"detail"`               //Tasks_detail
	Assigned_by   string   `db:"assigned_by" json:"assigned_by"`     //Personnels
	Assignment_id int      `db:"assignment_id" json:"assignment_id"` //Tasks_assignment
	File          string   `db:"file" json:"-"`                      //Tasks_attachments
	Uploaded_at   string   `db:"uploaded_at" json:"uploaded_at"`     //Tasks_attachments
	Files         []string `json:"files"`
}

func GetTaskEvidence(db *sqlx.DB, task_id int) ([]TaskEvidenceInfo, error) {

	var taskevidenceinfo []TaskEvidenceInfo
	query := `
        SELECT tas.task_id, tas.assignment_id, tat.file, tat.Uploaded_at, td.title, td.detail, p.first_name || ' ' || p.last_name AS assigned_by
        FROM "Tasks_assignment" tas
        LEFT JOIN "Tasks_detail" td ON tas.task_id = td.task_id
        LEFT JOIN "Personnels" p ON td.assigned_by = p.personnel_id
        LEFT JOIN "Tasks_attachments" tat ON tas.assignment_id = tat.assignment_id
		WHERE tas.task_id = $1
    `
	err := db.Select(&taskevidenceinfo, query, task_id)
	if err != nil {
		return nil, err
	}

	for i := range taskevidenceinfo {
		files := strings.Split(taskevidenceinfo[i].File, ",")
		for j := range files {
			files[j] = config.APIURL + strings.TrimSpace(files[j])
		}
		taskevidenceinfo[i].Files = files
	}

	return taskevidenceinfo, nil
}
