package main

import (
	"fmt"
	"html"
	"net/http"
	"net/url"
	"io"
	"encoding/json"
	"errors"
)

type Order struct {
	userName string
	userInf  string
}

func (order *Order) sendNotification() error {
	telegramMessageURL := fmt.Sprintf(globalConfig.TelegramURLTemplate, globalConfig.TelegramBotTokken)
	
	data := url.Values{
		"chat_id": {globalConfig.TelegramChatID},
		"text": {fmt.Sprintf(globalConfig.OrderMessageTemplate, html.EscapeString(order.userName), html.EscapeString(order.userInf))},
		"parse_mode": {"HTML"},
	}

	resp, err := http.PostForm(telegramMessageURL, data)
	if err != nil {
		return err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	answer := make(map[string]bool)
	json.Unmarshal(bodyBytes, &answer)

	if !answer["ok"] {
		return errors.New("telegram notify error")
	}

	return nil
}