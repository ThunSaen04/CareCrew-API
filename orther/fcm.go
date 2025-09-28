package orther

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"firebase.google.com/go/v4/messaging"
	fcm "github.com/appleboy/go-fcm"
	"github.com/jmoiron/sqlx"
	"github.com/project/carecrew/config"
)

type SendNotiInfo struct {
	Task_id int    `db:"task_id" json:"task_id"`
	Detail  string `json:"detail"`
	Title   string `json:"title"`
	Body    string `json:"body"`
}

func saveNotification(db *sqlx.DB, personnelID int, title, body string, data map[string]string) error {
	dataJSON := "{}"
	if len(data) > 0 {
		b, err := json.Marshal(data)
		if err != nil {
			return err
		}
		dataJSON = string(b)
	}

	var count int
	err := db.Get(&count, `
		SELECT COUNT(1) 
		FROM "Notifications" 
		WHERE personnel_id=$1 AND title=$2 AND body=$3 AND data=$4
	`, personnelID, title, body, dataJSON)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Printf("[System] Notification ซ้ำ ไม่บันทึก: personnel_id=%d title=%s", personnelID, title)
		return nil
	}

	_, err = db.Exec(`
		INSERT INTO "Notifications" (personnel_id, title, body, data)
		VALUES ($1, $2, $3, $4)
	`, personnelID, title, body, dataJSON)
	return err
}

func SendNotificationToAll(db *sqlx.DB, title, body string, taskid int) error {
	ctx := context.Background()
	client, err := fcm.NewClient(ctx, fcm.WithCredentialsFile(config.FullFCMPath))
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
		return nil
	}

	msg := &messaging.MulticastMessage{
		Notification: &messaging.Notification{Title: title, Body: body},
		Data: map[string]string{
			"task_id": fmt.Sprint(taskid),
			"title":   title,
			"body":    body,
		},
		Tokens:  tokens,
		Android: &messaging.AndroidConfig{Priority: "high"},
	}

	resp, err := client.SendMulticast(ctx, msg)
	if err != nil {
		return err
	}

	log.Print("[System] Success:", resp.SuccessCount, " Failure:", resp.FailureCount)

	var personnelIDs []int
	err = db.Select(&personnelIDs, `SELECT personnel_id FROM "Personnels"`)
	if err != nil {
		log.Print("[Error] ไม่สามารถเรียกข้อมูล personnel_id ได้:", err)
		return err
	}

	for _, personnelID := range personnelIDs {
		saveNotification(db, personnelID, title, body, map[string]string{
			"task_id": fmt.Sprint(taskid),
		})
	}

	return nil
}

func SendNotiSuccessToPerInTask(db *sqlx.DB, sendnotiinfo *SendNotiInfo) error {
	ctx := context.Background()
	client, err := fcm.NewClient(ctx, fcm.WithCredentialsFile(config.FullFCMPath))
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
		return nil
	}

	msg := &messaging.MulticastMessage{
		Notification: &messaging.Notification{Title: sendnotiinfo.Title, Body: sendnotiinfo.Body},
		Data: map[string]string{
			"task_id": fmt.Sprint(sendnotiinfo.Task_id),
			"detail":  fmt.Sprint(sendnotiinfo.Detail),
			"title":   sendnotiinfo.Title,
			"body":    sendnotiinfo.Body,
		},
		Tokens:  tokens,
		Android: &messaging.AndroidConfig{Priority: "high"},
	}

	resp, err := client.SendMulticast(ctx, msg)
	if err != nil {
		return err
	}

	log.Print("[System] Success:", resp.SuccessCount, " Failure:", resp.FailureCount)

	var personnelIDs []int
	err = db.Select(&personnelIDs, `
		SELECT personnel_id 
		FROM "Tasks_assignment" 
		WHERE task_id = $1
	`, sendnotiinfo.Task_id)
	if err != nil {
		log.Print("[Error] ไม่สามารถเรียกข้อมูล personnel_id ได้:", err)
		return err
	}

	for _, personnelID := range personnelIDs {
		saveNotification(db, personnelID, sendnotiinfo.Title, sendnotiinfo.Body, map[string]string{
			"task_id": fmt.Sprint(sendnotiinfo.Task_id),
			"detail":  fmt.Sprint(sendnotiinfo.Detail),
		})
	}

	return nil
}

func SendNotiSuccessToAssignor(db *sqlx.DB, sendnotiinfo *SendNotiInfo) error {
	ctx := context.Background()
	client, err := fcm.NewClient(ctx, fcm.WithCredentialsFile(config.FullFCMPath))
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
		return nil
	}

	msg := &messaging.MulticastMessage{
		Notification: &messaging.Notification{Title: sendnotiinfo.Title, Body: sendnotiinfo.Body},
		Data: map[string]string{
			"task_id": fmt.Sprint(sendnotiinfo.Task_id),
			"detail":  fmt.Sprint(sendnotiinfo.Detail),
			"title":   sendnotiinfo.Title,
			"body":    sendnotiinfo.Body,
		},
		Tokens:  tokens,
		Android: &messaging.AndroidConfig{Priority: "high"},
	}

	resp, err := client.SendMulticast(ctx, msg)
	if err != nil {
		return err
	}

	log.Print("[System] Success:", resp.SuccessCount, " Failure:", resp.FailureCount)

	var personnelIDs []int
	err = db.Select(&personnelIDs, `
		SELECT personnel_id 
		FROM "Personnels" 
		WHERE role_type_id = 1
	`)
	if err != nil {
		log.Print("[Error] ไม่สามารถเรียกข้อมูล personnel_id ได้:", err)
		return err
	}

	for _, personnelID := range personnelIDs {
		saveNotification(db, personnelID, sendnotiinfo.Title, sendnotiinfo.Body, map[string]string{
			"task_id": fmt.Sprint(sendnotiinfo.Task_id),
			"detail":  fmt.Sprint(sendnotiinfo.Detail),
		})
	}

	return nil
}
