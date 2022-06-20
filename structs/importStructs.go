package structs

import (
	"time"
)

type Pattern struct {
	Regex  string `json:"regex"`
	Format string `json:"pattern"`
}

type MonthData struct {
	Hashes    []string
	StartTime time.Time
	EndTime   time.Time
	Data      []DataPoint
}

type DataPoint struct {
	Start time.Time
	Path  string `hash:"-"`
}
