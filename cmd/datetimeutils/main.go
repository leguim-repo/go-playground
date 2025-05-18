package main

import (
	"fmt"
	"go-playground/pkg/datetimeutils"
)

func main() {

	fmt.Println("Example of use datetimeutils module:")
	// Using Now()
	currentTimeUTC := datetimeutils.Now()
	fmt.Println("Current time in UTC:", currentTimeUTC)

	// Using CreateTimestamp()
	// Let's use the time obtained from Now() as an example input
	formattedTimeStamp, err := datetimeutils.CreateTimeStamp(currentTimeUTC)
	if err != nil {
		fmt.Println("Error creating timestamp:", err)
	} else {
		fmt.Println("Formatted timestamp in Europe/Berlin:", formattedTimeStamp)
	}

	// Using CreatePartitionStamp()
	// Let's use the time obtained from Now() as an example input
	formattedPartitionStamp, err := datetimeutils.CreatePartitionStamp(currentTimeUTC)
	if err != nil {
		fmt.Println("Error creating timestamp:", err)
	} else {
		fmt.Println("Formatted partition stamp in Europe/Berlin:", formattedPartitionStamp)
	}

	// Using UnixTimestampWithMilliseconds()
	unixMillis := datetimeutils.UnixTimestampWithMilliseconds()
	fmt.Println("Current Unix timestamp in milliseconds:", unixMillis)

	// Using DateAndTimeFromTod()
	// Example timestamp string (milliseconds)
	timeFromUnixMillis, err := datetimeutils.ConvertUnixToDateTime(unixMillis)
	if err != nil {
		fmt.Println("Error converting unix timestamp string:", err)
	} else {
		fmt.Println("Time from Unix milliseconds string:", timeFromUnixMillis)
	}

	// Using ConvertToUnixTimestamp()
	// Example date string matching the timestampLayout
	exampleDateStr := datetimeutils.ConvertTimeStampToDateStr(formattedTimeStamp) // Example: March 15, 2023 10:00:00.500000 UTC
	unixSeconds, err := datetimeutils.ConvertToUnixTimestamp(exampleDateStr)
	if err != nil {
		fmt.Println("Error converting date string to unix timestamp:", err)
	} else {
		fmt.Println("Unix timestamp (seconds) from date string:", unixSeconds)
	}
}
