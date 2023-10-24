package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

type DayCost struct {
	Day      *string `json:"day"`
	Cost     float64 `json:"cost"`
	Resource string  `json:"resource"`
}

type CostExplorer struct {
	cli *costexplorer.CostExplorer
}

func NewCostExplorer() (*CostExplorer, error) {
	s, err := session.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}
	return &CostExplorer{
		cli: costexplorer.New(s),
	}, nil
}

func (c *CostExplorer) GetCostAndUsage(start time.Time, end time.Time) (*costexplorer.GetCostAndUsageOutput, error) {
	graudarity := aws.String("DAILY")
	timePeriod := &costexplorer.DateInterval{
		Start: aws.String(start.Format("2006-01-02")),
		End:   aws.String(end.Format("2006-01-02")),
	}
	groupBy := []*costexplorer.GroupDefinition{
		{
			// Key:  aws.String("Hoge"),
			// Type: aws.String("TAG"),
			Key:  aws.String("SERVICE"),
			Type: aws.String("DIMENSION"),
		},
	}
	metrics := []*string{
		aws.String("BlendedCost"),
	}

	input := &costexplorer.GetCostAndUsageInput{
		Granularity: graudarity,
		TimePeriod:  timePeriod,
		GroupBy:     groupBy,
		Metrics:     metrics,
	}

	o, err := c.cli.GetCostAndUsage(input)
	if err != nil {
		return nil, fmt.Errorf("failed to get cost and usage: %v", err)
	}
	return o, nil
}

func FormatCost(cost *costexplorer.GetCostAndUsageOutput) ([]*DayCost, error) {
	res := cost.ResultsByTime
	var days []*DayCost
	for _, c := range res {
		for _, s := range c.Groups {
			f, err := strconv.ParseFloat(*s.Metrics["BlendedCost"].Amount, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse float: %v", err)
			}
			days = append(days, &DayCost{
				Day:      c.TimePeriod.Start,
				Cost:     f,
				Resource: *s.Keys[0],
			})
		}
	}
	return days, nil
}

func main() {
	// cost explorer

	// post to slack

	c, err := NewCostExplorer()
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	res, err := c.GetCostAndUsage(start, end)
	if err != nil {
		log.Fatal(err)
	}

	days, err := FormatCost(res)
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.MarshalIndent(days, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create("cost.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.Write(b)

}

// func FetchCost() []Cost {
// 	// fetch cost from aws
// 	return []Cost{}
// }
//
// func PostToSlack(costs []Cost) {
// 	// post to slack
// }
