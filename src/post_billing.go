package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/tabwriter"
	"time"
)

func PostBilling() {
	cfg, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	c, err := NewCostExplorer()
	if err != nil {
		log.Fatal(err)
	}

	client := http.Client{}
	slack, err := NewSlackClient(&client, cfg.WebHookUrl)
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)       // 月初
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC) // 今日
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

	dd := fmt.Sprintf(
		"%s ~ %s",
		start.Format("2006-01-02"),
		end.Format("2006-01-02"))

	message = fmt.Sprintf("%s\n```%s```", dd, message)

	if err := slack.SendMessage(message); err != nil {
		log.Fatal(err)
	}
}

func FormatMessage(c ResourceCost) (string, error) {
	var buf bytes.Buffer
	tWriter := tabwriter.NewWriter(&buf, 0, 0, 0, ' ', tabwriter.Debug)
	fmt.Fprintln(tWriter, "Service", "\t", "Cost($)")
	fmt.Fprintln(tWriter, "==========", "\t", "==========")
	for _, a := range c.ToArray() {
		fmt.Fprintln(tWriter, a[0], "\t", a[1])
	}
	total, err := c.Total()
	if err != nil {
		return "", fmt.Errorf("failed to get total cost: %w", err)
	}
	fmt.Fprintln(tWriter, "----------", "\t", "----------")
	fmt.Fprintln(tWriter, "Total", "\t", strconv.FormatFloat(total, 'f', 2, 64))
	tWriter.Flush()
	return buf.String(), nil
}
