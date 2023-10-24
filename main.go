package main

import (
	"log"
	"time"
)

func main() {
	cfg, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	c, err := NewCostExplorer()
	if err != nil {
		log.Fatal(err)
	}

	slack := NewSlackClient(cfg.WebHookUrl)

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	res, err := c.GetCostAndUsage(start, end)
	if err != nil {
		log.Fatal(err)
	}

	days, err := FormatFlatCost(res)
	if err != nil {
		log.Fatal(err)
	}

	fc, err := AggrigateCost(days)
	if err != nil {
		log.Fatal(err)
	}

	message, err := FormatMessage(fc)
	if err != nil {
		log.Fatal(err)
	}

	if err := slack.SendMessage(message, cfg.SlackChannelBilling); err != nil {
		log.Fatal(err)
	}
}
