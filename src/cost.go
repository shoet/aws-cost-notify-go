package main

import (
	"fmt"
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

type Costs []*DayCost

type ResourceCost map[string]float64

func (c Costs) AggrigateResource() ResourceCost {
	tmp := make(map[string]float64)
	for _, v := range c {
		tmp[v.Resource] += v.Cost
	}
	return tmp
}

func (r ResourceCost) ToArray() [][]string {
	/*
		example: {{"EC", 100.00}, {"S3", 200.00}}
	*/
	arr := make([][]string, 0, len(r))
	for k, v := range r {
		arr = append(arr, []string{k, strconv.FormatFloat(v, 'f', 2, 64)})
	}
	return arr
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

func (c *CostExplorer) GetCostAndUsage(
	start time.Time, end time.Time,
) (*costexplorer.GetCostAndUsageOutput, error) {
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

func FormatFlatCost(cost *costexplorer.GetCostAndUsageOutput) ([]*DayCost, error) {
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

func AggrigateCost(c Costs) (ResourceCost, error) {
	return c.AggrigateResource(), nil
}
