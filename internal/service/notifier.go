package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type NotifierService struct {
	TelegramToken string
	TelegramChatID string
	SlackWebhookURL string
	UseTelegram bool
	UseSlack bool
}

func NewNotifierService(token, chatID, slackURL string, useTelegram, useSlack bool) *NotifierService {
	return &NotifierService{
		TelegramToken: token,
		TelegramChatID: chatID,
		SlackWebhookURL: slackURL,
		UseTelegram: useTelegram,
		UseSlack: useSlack,
	}
}

func (n *NotifierService) Send(repo string, success bool) {
	status := "✅ Pipeline succeeded"
	if !success {
		status = "❌ Pipeline failed"
	}

	message := fmt.Sprintf("%s\nRepository: %s", status, repo)

	if n.UseTelegram {
		n.sendTelegram(message)
	}
	if n.UseSlack {
		n.sendSlack(message)
	}
}

func (n *NotifierService) sendTelegram(text string) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", n.TelegramToken)

	body := map[string]string{
		"chat_id": n.TelegramChatID,
		"text":    text,
	}
	jsonBody, _ := json.Marshal(body)

	_, _ = http.Post(url, "application/json", bytes.NewReader(jsonBody))
}

func (n *NotifierService) sendSlack(text string) {
	body := map[string]string{"text": text}
	jsonBody, _ := json.Marshal(body)

	_, _ = http.Post(n.SlackWebhookURL, "application/json", bytes.NewReader(jsonBody))
}
