package orther

import (
	"context"
	"log"

	"firebase.google.com/go/v4/messaging"
	fcm "github.com/appleboy/go-fcm"
	"github.com/jmoiron/sqlx"
	"github.com/project/carecrew/config"
)

/* func SendNotificationToAll(db *sqlx.DB, title, body string) {
	var tokens []string
	err := db.Select(&tokens, `SELECT token FROM "FCM_Tokens"`)
	if err != nil {
		log.Print("[Error] ไม่สามารถเรียก Token ได้:", err)
	}

	for _, token := range tokens {
		message := map[string]interface{}{
			"to": token,
			"notification": map[string]string{
				"title": title,
				"body":  body,
			},
		}
	}
} */

func SendNotificationToAll(db *sqlx.DB, title, body string) error {
	ctx := context.Background()
	client, err := fcm.NewClient(
		ctx,
		fcm.WithCredentialsFile(config.FullFCMPath),
	)
	if err != nil {
		return err
	}

	var tokens []string
	err = db.Select(&tokens, `SELECT token FROM "FCM_Tokens"`)
	if err != nil {
		log.Print("[Error] ไม่สามารถเรียกข้อมูล Token ได้:", err)
		return err
	}
	if len(tokens) == 0 {
		log.Print("[Warning] ไม่พบ Token ใน DB")
		return err
	}
	//token := "cMtcUYSqTZKogwbS4vEjM_:APA91bEAA1PhVC6ZwKkNvGqBaGuzahUf7q4zLDZNVsNDus-PVLneLBrykiExwZyyZdWC-lgplHtfTzo_orryGrUX_VCHj2ll20icnlae6ZV7LkHzGYLj4oc"
	msg := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Tokens: tokens,
	}
	resp, err := client.SendMulticast(ctx, msg)
	if err != nil {
		return err
	}
	log.Print("[System] ", "Success:", resp.SuccessCount, "Failure:", resp.FailureCount)

	return nil
}

type SendNotiInfo struct {
	Task_id int    `db:"task_id" json:"task_id"`
	Title   string `json:"title"`
	Body    string `json:"body"`
}

func SendNotiSuccessToPerInTask(db *sqlx.DB, sendnotiinfo *SendNotiInfo) error {
	ctx := context.Background()
	client, err := fcm.NewClient(
		ctx,
		fcm.WithCredentialsFile(config.FullFCMPath),
	)
	if err != nil {
		return err
	}

	var tokens []string
	err = db.Select(&tokens, `
        SELECT DISTINCT t.token
        FROM "Tasks_assignment" a
        JOIN "FCM_Tokens" t ON t.personnel_id = a.personnel_id
        WHERE a.task_id = $1
    `, sendnotiinfo.Task_id)
	if err != nil {
		return err
	}
	if len(tokens) == 0 {
		log.Print("[Warning] ไม่พบ Token")
		return err
	}
	msg := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: sendnotiinfo.Title,
			Body:  sendnotiinfo.Body,
		},
		Tokens: tokens,
	}
	resp, err := client.SendMulticast(ctx, msg)
	if err != nil {
		return err
	}
	log.Print("Success:", resp.SuccessCount, "Failure:", resp.FailureCount)
	return nil
}

func SendNotiSuccessToAssignor(db *sqlx.DB, sendnotiinfo *SendNotiInfo) error {
	ctx := context.Background()
	client, err := fcm.NewClient(
		ctx,
		fcm.WithCredentialsFile(config.FullFCMPath),
	)
	if err != nil {
		return err
	}

	var tokens []string
	err = db.Select(&tokens, `
		SELECT DISTINCT t.token
		FROM "Personnels" p
		JOIN "FCM_Tokens" t ON t.personnel_id = p.personnel_id
		WHERE p.role_type_id = 1;
    `)
	if err != nil {
		return err
	}
	if len(tokens) == 0 {
		log.Print("[Warning] ไม่พบ Token")
		return err
	}
	msg := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: sendnotiinfo.Title,
			Body:  sendnotiinfo.Body,
		},
		Tokens: tokens,
	}
	resp, err := client.SendMulticast(ctx, msg)
	if err != nil {
		return err
	}
	log.Print("Success:", resp.SuccessCount, "Failure:", resp.FailureCount)
	return nil
}
