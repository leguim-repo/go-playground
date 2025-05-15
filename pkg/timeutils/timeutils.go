package timeutils

import (
	"strconv"
	"time"
)

// GetUnixTimestampWithMilliseconds return a Unix timestamp with milliseconds
func GetUnixTimestampWithMilliseconds() string {
	timestamp := time.Now().UnixMilli()
	return strconv.FormatInt(timestamp, 10)
}
