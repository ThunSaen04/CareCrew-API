package orther

import (
	"context"
	"fmt"
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
		log.Println("[Error] ไม่สามารถเรียกข้อมูล Token ได้:", err)
		return err
	}
	if len(tokens) == 0 {
		log.Println("[Warning] ไม่พบ Token ใน DB")
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
	fmt.Println("Success:", resp.SuccessCount, "Failure:", resp.FailureCount)

	return nil
}
