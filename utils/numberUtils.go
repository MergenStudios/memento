package utils

func PositiveTimestmap(timestamp int64) int64 {
	if timestamp < 0 {
		return -1
	} else {
		return timestamp
	}
}