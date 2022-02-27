package utils

import "memento/structs"

func GetDatapointsLen(dataPoints map[string][]structs.DataPoint) int64 {
	var count int64 = 0
	for _, month := range dataPoints {
		for _, _ = range month {
			count++
		}
	}
	return count
}
