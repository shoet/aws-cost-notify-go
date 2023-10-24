package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

type Cost struct {
	Tag      string
	Cost     float64
	Resource string
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

func (c *CostExplorer) GetCostAndUsage(start string, end string) (*costexplorer.GetCostAndUsageOutput, error) {
	graudarity := aws.String("DAILY")
	timePeriod := &costexplorer.DateInterval{
		Start: aws.String("2023-10-22"),
		End:   aws.String("2023-10-23"),
	}
	groupBy := []*costexplorer.GroupDefinition{
		{
			Key:  aws.String("SERVICE"),
			Type: aws.String("DIMENSION"),
		},
	}
	metrics := []*string{
		aws.String("UnblendedCost"),
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

func main() {
	// cost explorer

	// post to slack

	c, err := NewCostExplorer()
	if err != nil {
		log.Fatal(err)
	}

	res, err := c.GetCostAndUsage("2023-10-22", "2023-10-23")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)

}

func FetchCost() []Cost {
	// fetch cost from aws
	return []Cost{}
}

func PostToSlack(costs []Cost) {
	// post to slack
}

func FormatCost(costs []Cost) string {
	// format cost
	// texts := []string{}
	// for _, c := range costs {
	// 	texts = append(texts, "")
	// }
	return ""
}
