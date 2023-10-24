package main

import (
	"github.com/caarlos0/env/v9"
)

type Config struct {
	WebHookUrl          string `env:"SLACK_WEBHOOK_URL,required"`
	SlackChannelBilling string `env:"SLACK_CHANNEL_BILLING,required"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
