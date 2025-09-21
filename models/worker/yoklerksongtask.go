package worker

import (
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
)

type YokLerkSongTaskInfo struct {
	Personnel_id int `db:"personnel_id" json:"personnel_id"`
	Task_id      int `db:"task_id" json:"task_id"`
}

type fileinatt struct {
	File string `db:"file" json:"file"`
}

type s struct {
	Completed bool `db:"completed" json:"completed"`
}

func Yoklerksongtask(db *sqlx.DB, yoklerksongtaskinfo *YokLerkSongTaskInfo) error {

	tranX, err := db.Beginx()
	if err != nil {
		return err
	}

	defer tranX.Rollback()

	//เช็คว่ามีงานนั้นจริงอ๊ะป่าว
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1 FROM "Tasks_assignment"
			WHERE task_id = $1 AND personnel_id = $2 AND submit = true
		)
	`
	err = tranX.Get(&exists, query, yoklerksongtaskinfo.Task_id, yoklerksongtaskinfo.Personnel_id)
	if err != nil {
		return err
	}
	if exists { //ถ้ามีงานจริง
		//เช็คว่างานสถานะมันเป็นเสร็จสิ้นไหม ป้องกันตอนที่งานเสร็จแล้วผู้เล่นยกเลิก (กันไว้ก่อน)
		var checktask s
		query = `
			SELECT t.completed
			FROM "Tasks" t
			WHERE t.task_id = $1
		`
		err = tranX.Get(&checktask, query, yoklerksongtaskinfo.Task_id)
		if err != nil {
			log.Printf("[Warning] ไม่พบงานหมายเลข %d", yoklerksongtaskinfo.Task_id)
			return err
		}

		if !checktask.Completed {
			// เปลี่ยนสถานะการส่งงานกลับเป็น false
			var assignment_id int
			err = tranX.QueryRow(`
				UPDATE "Tasks_assignment"
				SET submit = false, submit_at = NULL
				WHERE personnel_id = $1 AND task_id = $2
				RETURNING assignment_id
			`, yoklerksongtaskinfo.Personnel_id, yoklerksongtaskinfo.Task_id).Scan(&assignment_id)
			if err != nil {
				return err
			}

			// เปลี่ยนสถานะงานกลับเป็น กำลังดำเนินการ
			_, err = tranX.Exec(`
				UPDATE "Tasks"
				SET status_type_id = 2, updated_at = NOW()
				WHERE task_id = $1
			`, yoklerksongtaskinfo.Task_id)
			if err != nil {
				return err
			}

			// เช็คว่ามีการส่งแล้วข้อมูลจริงใหม่
			query = `
				SELECT EXISTS (
					SELECT 1 FROM "Tasks_attachments"
					WHERE assignment_id = $1
				)
			`
			var exists2 bool
			err = tranX.Get(&exists2, query, assignment_id)
			if err != nil {
				return err
			}
			if exists2 { // ถ้ามีจริง ให้ลบออก (ลบไฟล์รูปออกด้วย)
				var getfilepath fileinatt
				query = `
					SELECT tat.file
					FROM "Tasks_attachments" tat
					WHERE assignment_id = $1
				`
				err = tranX.Get(&getfilepath, query, assignment_id)
				if err != nil {
					return err
				}

				_, err = tranX.Exec(`
					DELETE FROM "Tasks_attachments"
					WHERE assignment_id = $1
				`, assignment_id)
				if err != nil {
					return err
				}

				// path := strings.Split(getfilepath.File, ",")
				// for _, p := range path {
				// 	p = strings.TrimSpace(p)
				// 	if p != "" {
				// 		fullpath := config.BasePath + p

				// 		err := os.Remove(fullpath)
				// 		if err != nil {
				// 			if os.IsNotExist(err) {
				// 				log.Printf("ไม่พบไฟล์: %s", fullpath)
				// 			} else {
				// 				return fmt.Errorf("ลบไฟล์ไม่สำเร็จ %s: %w", fullpath, err)
				// 			}
				// 		} else {
				// 			log.Printf("ลบไฟล์สำเร็จ: %s", fullpath)
				// 		}
				// 	}
				// }
			} else {
				return errors.New("ไม่พบข้อมูลการส่งงาน")
			}
		}
	} else {
		return errors.New("ไม่พบงานการส่งงานนี้")
	}

	_, err = tranX.Exec(`
			INSERT INTO "Worker_logs" (personnel_id, task_id, detail)
			VALUES ($1, $2, $3)
		`, yoklerksongtaskinfo.Personnel_id, yoklerksongtaskinfo.Task_id, "ยกเลิกส่งงาน")
	if err != nil {
		return err
	}

	err = tranX.Commit()
	if err != nil {
		return err
	}
	return nil
}
