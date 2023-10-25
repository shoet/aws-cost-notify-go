package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client interface {
	Do(*http.Request) (*http.Response, error)
}

type SlackClient struct {
	client     Client
	webHookUrl string
}

func NewSlackClient(client Client, webHookUrl string) (*SlackClient, error) {
	return &SlackClient{
		client:     client,
		webHookUrl: webHookUrl,
	}, nil
}

func (s *SlackClient) SendMessage(message string, channelId string) error {
	post := struct {
		Text string `json:"text"`
	}{
		Text: message,
	}
	b, err := json.Marshal(post)
	req, err := http.NewRequest("POST", s.webHookUrl, bytes.NewBuffer(b))
	_, err = s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	return nil
}
