package main

import (
	"net/http"
	"testing"
)

func Test_SlackClient_SendMessage(t *testing.T) {
	client := http.Client{}
	cfg, err := NewConfig()
	if err != nil {
		t.Fatalf("failed to create config: %v", err)
	}
	slack, err := NewSlackClient(&client, cfg.WebHookUrl)
	if err != nil {
		t.Fatalf("failed to create slack client: %v", err)
	}
	slack.SendMessage("test")
}
