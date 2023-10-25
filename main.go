package main

import (
	"bytes"
	"fmt"
	"log"
	"text/tabwriter"
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

	fmt.Println(message)
	_ = slack

	// if err := slack.SendMessage(message, cfg.SlackChannelBilling); err != nil {
	// 	log.Fatal(err)
	// }
}

func FormatMessage(c ResourceCost) (string, error) {
	var buf bytes.Buffer
	tWriter := tabwriter.NewWriter(&buf, 0, 0, 0, ' ', tabwriter.Debug)
	fmt.Fprintln(tWriter, "Service", "\t", "Cost($)")
	for _, a := range c.ToArray() {
		fmt.Fprintln(tWriter, a[0], "\t", a[1])
	}
	tWriter.Flush()
	return buf.String(), nil
}
