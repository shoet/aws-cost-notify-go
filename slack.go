package main

type SlackClient struct {
	WebHookUrl string
}

func NewSlackClient(webHookUrl string) *SlackClient {
	return &SlackClient{WebHookUrl: webHookUrl}
}

func (s *SlackClient) SendMessage(message string, channelId string) error {
	return nil
}
