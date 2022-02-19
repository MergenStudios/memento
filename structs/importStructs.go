package structs

import (
	"time"
)

type TypeEnum struct {
	TrueName      string   `json:"true_name"`
	Extensions    []string `json:"extensions"`
	Dated         string   `json:"Dated"`
	DetermineTime string   `json:"determine_time"`
}

type MonthData struct {
	StartTime time.Time
	EndTime   time.Time
	Data      []DataPoint
}

type DataPoint struct {
	Dated string
	Start time.Time
	End   time.Time
	Type  string
	Path  string
}
