package main

import "time"

type Clocker interface {
	Now() time.Time
}

type RealClocker struct{}

func (c *RealClocker) Now() time.Time {
	return time.Now()
}

func GetCurrentMonthPeriod(c Clocker) (time.Time, time.Time) {
	now := c.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC) // 月初
	end := start.AddDate(0, 1, -1)                                       // 月末
	return start, end
}
