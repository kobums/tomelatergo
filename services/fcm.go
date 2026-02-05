package services

import (
	"context"
	"tomelater/global"
	"tomelater/global/log"
	"tomelater/models"
	"tomelater/models/alarm"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

func Fcm() {
	log.Info().Str("service", "FCM").Msg("Start Service")

	go func() {
		ch := global.GetFcm()

		ctx := context.Background()
		opt := option.WithCredentialsFile("./fcm.json")
		app, err := firebase.NewApp(ctx, nil, opt)
		if err != nil {
			log.Printf("error initializing app: %v\n", err)
			return
		}

		client, err := app.Messaging(context.Background())
		if err != nil {
			log.Printf("Error initializing Firebase Cloud Messaging client: %v\n", err)
		}

		for item := range ch {
			conn := models.NewConnection()

			tokenManager := models.NewTokenManager(conn)
			alarmManager := models.NewAlarmManager(conn)

			for _, user := range item.Target {
				tokens := tokenManager.Find([]interface{}{
					models.Where{Column: "user", Value: user, Compare: "="},
				})

				status := alarm.StatusFail

				for _, token := range tokens {
					message := &messaging.Message{
						Data:         item.Message,
						Notification: &messaging.Notification{Title: item.Title, Body: item.Message["message"]},
						Token:        token.Token,
					}

					// 메시지 전송
					response, err := client.Send(context.Background(), message)
					if err != nil {
						log.Printf("Error sending message: %v\n", err)
					} else {
						log.Println("Successfully sent message:", response)
						status = alarm.StatusSuccess
					}
				}

				alarmItem := models.Alarm{
					Title:   item.Title,
					Content: item.Message["message"],
					Type:    alarm.Type(item.Type),
					Status:  status,
					User:    user,
				}

				alarmManager.Insert(&alarmItem)
			}

			conn.Close()
		}
	}()
}
