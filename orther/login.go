package orther

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type loginInfo struct {
	PersonnelID int    `db:"personnel_id" json:"personnel_id"` //Personnels
	Password    string `db:"password" json:"-"`                //User
	Role        string `db:"role_type_id" json:"role_type_id"` //Role_types
}

func Auth(db *sqlx.DB, PersonnelID int, Password string) (*loginInfo, error) {
	var login loginInfo
	query := `
		SELECT u.personnel_id, u.password, p.role_type_id
		FROM "Users" u
		INNER JOIN "Personnels" p ON u.personnel_id = p.personnel_id
		WHERE u.personnel_id = $1
	`

	err := db.Get(&login, query, PersonnelID)
	if err != nil {
		return nil, err
	}

	// if Password != login.Password {
	// 	return nil, errors.New("ฟหกาฟยนไก่ยฟนไก่")
	// }

	CheckHash := CheckPasswordHash(Password, login.Password)

	if CheckHash != true {
		// Println("PersonnelID = ", PersonnelID)
		// Println("Password = ", Password)
		// Println("login.Password = ", login.Password)
		// Println("CheckHash = ", CheckHash)
		return nil, errors.New("รหัสผิด")
	}

	_, err = db.Exec(
		`UPDATE "Users"
		SET last_login = NOW()
		WHERE personnel_id = $1`,
		PersonnelID,
	)
	if err != nil {
		return nil, err
	}

	// Println("PersonnelID = ", PersonnelID)
	// Println("Password = ", Password)
	// Println("login.Password = ", login.Password)
	// Println("CheckHash = ", CheckHash)
	return &login, nil
}

func AuthV2(db *sqlx.DB, PersonnelID int, Password string, Token string) (*loginInfo, error) {

	tranX, err := db.Beginx()

	var login loginInfo
	query := `
		SELECT u.personnel_id, u.password, p.role_type_id
		FROM "Users" u
		INNER JOIN "Personnels" p ON u.personnel_id = p.personnel_id
		WHERE u.personnel_id = $1
	`

	defer tranX.Rollback()

	err = tranX.Get(&login, query, PersonnelID)
	if err != nil {
		return nil, errors.New("หมายเลขผู้ใช้งานหรือรหัสผ่านไม่ถูกต้อง")
	}

	CheckHash := CheckPasswordHash(Password, login.Password)
	if CheckHash != true {
		return nil, errors.New("หมายเลขผู้ใช้งานหรือรหัสผ่านไม่ถูกต้อง")
	}

	_, err = tranX.Exec(
		`UPDATE "Users"
		SET last_login = NOW()
		WHERE personnel_id = $1`,
		PersonnelID)
	if err != nil {
		return nil, err
	}

	if Token != "" {
		_, err = tranX.Exec(`
            INSERT INTO "FCM_Tokens" (personnel_id, token)
            VALUES ($1, $2)
            ON CONFLICT (token) DO UPDATE 
            SET personnel_id = $1, updated_at = NOW()
        `, PersonnelID, Token)
		if err != nil {
			return nil, err
		}
	}

	err = tranX.Commit()
	if err != nil {
		return nil, err
	}

	//log.Println(PersonnelID, Password, Token)

	return &login, nil
}
