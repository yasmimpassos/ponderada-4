package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type NotificationService struct{}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

type fcmNotification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type fcmPayload struct {
	To           string          `json:"to"`
	Notification fcmNotification `json:"notification"`
}

func (s *NotificationService) SendPush(fcmToken string, title string, body string) error {
	serverKey := os.Getenv("FCM_SERVER_KEY")
	if serverKey == "" || fcmToken == "" {
		return nil
	}

	payload := fcmPayload{
		To:           fcmToken,
		Notification: fcmNotification{Title: title, Body: body},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "https://fcm.googleapis.com/fcm/send", bytes.NewReader(payloadBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("key=%s", serverKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("FCM retornou status inesperado")
	}
	return nil
}

func (s *NotificationService) NotifyGroupMembers(tokens []string, title string, body string) {
	for _, token := range tokens {
		if token == "" {
			continue
		}
		_ = s.SendPush(token, title, body)
	}
}
