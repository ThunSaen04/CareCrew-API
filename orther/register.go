package orther

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
)

type RegisterUserInfo struct {
	Password   string `db:"password" json:"password"`         //User
	FirstName  string `db:"first_name" json:"first_name"`     //ersonnels
	LastName   string `db:"last_name" json:"last_name"`       //Personnels
	Phone      string `db:"phone" json:"phone"`               //Personnels
	Role_types int    `db:"role_type_id" json:"role_type_id"` //Role_types
}

func RegisterUser(db *sqlx.DB, newRegisterUserInfo *RegisterUserInfo) (int, error) {

	tranX, err := db.Beginx()

	var personnelID int
	hashPW, err := HashPassword(newRegisterUserInfo.Password)
	if err != nil {
		return 0, err
	}

	// println("Password = ", newRegisterUserInfo.Password)
	// println("HashPassword = ", hashPW)
	// println("FirstName = ", newRegisterUserInfo.FirstName)
	// println("LastName = ", newRegisterUserInfo.LastName)
	// println("Phone = ", newRegisterUserInfo.Phone)
	// println("Role_types = ", newRegisterUserInfo.Role_types)

	defer tranX.Rollback()

	err = tranX.QueryRow(
		`INSERT INTO "Personnels" (first_name, last_name, Phone, role_type_id)
		VALUES ($1, $2, $3, $4)
		RETURNING personnel_id`,
		newRegisterUserInfo.FirstName,
		newRegisterUserInfo.LastName,
		newRegisterUserInfo.Phone,
		newRegisterUserInfo.Role_types,
	).Scan(&personnelID)
	if err != nil {
		return 0, err
	}

	// println("PersonnelID = ", personnelID)

	_, err = tranX.Exec(
		`INSERT INTO "Users" (personnel_id, password)
		VALUES ($1, $2)`,
		personnelID,
		hashPW,
	)
	if err != nil {
		log.Error(err.Error())
		return 0, err
	}

	err = tranX.Commit()
	if err != nil {
		log.Error(err.Error())
		return 0, err
	}

	return personnelID, nil
}
