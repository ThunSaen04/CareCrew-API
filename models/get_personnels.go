package models

import (
	"github.com/jmoiron/sqlx"
)

type PersonnelsInfo struct {
	PersonnelID int    `db:"personnel_id" json:"personnel_id"` //Personnels
	FirstName   string `db:"first_name" json:"first_name"`     //Personnels
	LastName    string `db:"last_name" json:"last_name"`       //Personnels
	Phone       string `db:"phone" json:"phone"`               //Personnels
	RoleName    string `db:"name" json:"role_name"`            //Role_types
}

func GetPersonnelsInfo(db *sqlx.DB) ([]PersonnelsInfo, error) {
	personnels := []PersonnelsInfo{}
	query := `
        SELECT p.personnel_id, p.first_name, p.last_name, p.phone, r.name
        FROM "Personnels" p
        LEFT JOIN "Role_types" r ON p.role_type_id = r.role_type_id
		ORDER BY p.personnel_id DESC
    `

	err := db.Select(&personnels, query)
	return personnels, err
}
