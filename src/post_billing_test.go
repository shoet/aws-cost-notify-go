package main

import (
	"fmt"
	"testing"
	"time"
)

func Test_PostBilling(t *testing.T) {
	PostBilling()
}

type MockClocker struct {
	base time.Time
}

func (c *MockClocker) Now() time.Time {
	return time.Date(c.base.Year(), c.base.Month(), c.base.Day(), 0, 0, 0, 0, time.UTC)
}

func NewMockClocker(base time.Time) *MockClocker {
	return &MockClocker{
		base: base,
	}
}

func Test_GetCurrentMonthPeriod(t *testing.T) {
	base := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 100; i++ {
		c := NewMockClocker(base)
		start, end := GetCurrentMonthPeriod(c)
		fmt.Println(start, end)
		base = base.AddDate(0, 0, 1)
	}

}
