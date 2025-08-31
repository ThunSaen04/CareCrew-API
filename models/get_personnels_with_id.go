package models

import (
	"github.com/jmoiron/sqlx"
)

func GetPersonnelsInfo_With_ID(db *sqlx.DB, PersonnelID int) (*PersonnelsInfo, error) {
	var personnels PersonnelsInfo
	query := `
        SELECT p.personnel_id, p.first_name, p.last_name, p.phone, r.name
        FROM "Personnels" p
        LEFT JOIN "Role_types" r ON p.role_type_id = r.role_type_id
		WHERE p.personnel_id = $1
    `

	err := db.Get(&personnels, query, PersonnelID)
	if err != nil {
		return nil, err
	}

	return &personnels, nil
}
