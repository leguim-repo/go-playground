package datetimeutils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// GetUnixTimestampWithMilliseconds return a Unix timestamp with milliseconds
func GetUnixTimestampWithMilliseconds() string {
	timestamp := time.Now().UnixMilli()
	return strconv.FormatInt(timestamp, 10)
}

// Define the layouts corresponding to the Python formats
// Go's time formatting uses a reference time: Mon Jan 2 15:04:05 MST 2006
const (
	partitionLayout = "20060102"
	timestampLayout = "20060102150405.000000" // .000000 for microseconds
)

// Now gets the current time in UTC, removing timezone information.
// In Go, a time.Time object always has a location. Returning time.UTC()
// provides the time in the UTC location, which is the idiomatic equivalent
// of a timezone-naive UTC time in Python.
func Now() time.Time {
	// time.Now() gives the current local time.
	// Use time.UTC() to get the current time in UTC.
	return time.Now().UTC()
}

// CreateTimeStamp formats a given time.Time object into a string
// using the specified format layout and timezone (Europe/Berlin).
func CreateTimeStamp(fileCreationDate time.Time) (string, error) {
	// Load the Europe/Berlin timezone location
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		// Return the error if the location cannot be loaded
		return "", fmt.Errorf("failed to load timezone Europe/Berlin: %w", err)
	}

	// Convert the input time to the Berlin timezone
	timeInBerlin := fileCreationDate.In(loc)

	timeInBerlinValue := timeInBerlin.Format(timestampLayout)
	timeInBerlinValue = strings.Replace(timeInBerlinValue, ".", "", -1)

	// Format the time using the defined timestamp layout
	return timeInBerlinValue, nil
}

func CreatePartitionStamp(fileCreationDate time.Time) (string, error) { // Load the Europe/Berlin timezone location
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		// Return the error if the location cannot be loaded
		return "", fmt.Errorf("failed to load timezone Europe/Berlin: %w", err)
	}

	// Convert the input time to the Berlin timezone
	timeInBerlin := fileCreationDate.In(loc)

	// Format the time using the defined partition layout
	return timeInBerlin.Format(partitionLayout), nil
}

// UnixTimestampWithMilliseconds returns the current Unix timestamp in milliseconds.
func UnixTimestampWithMilliseconds() int64 {
	// time.Now().UTC() gets the current time in UTC.
	// UnixNano() gets the timestamp in nanoseconds.
	// Divide by 1e6 to convert nanoseconds to milliseconds.
	return time.Now().UTC().UnixNano() / int64(time.Millisecond)
}

// DateAndTimeFromTod converts a Unix timestamp string (in milliseconds)
// back to a time.Time object.
func DateAndTimeFromTod(unixTimestampStr string) (time.Time, error) {
	// Parse the string as a float64 to handle potential fractional milliseconds
	millisecondsFloat, err := strconv.ParseFloat(unixTimestampStr, 64)
	if err != nil {
		// Return zero time and error if parsing fails
		return time.Time{}, fmt.Errorf("failed to parse unix timestamp string: %w", err)
	}

	// Convert milliseconds to seconds (float)
	secondsFloat := millisecondsFloat / 1000.0

	// Separate integer seconds and fractional seconds
	seconds := int64(secondsFloat)
	nanoseconds := int64((secondsFloat - float64(seconds)) * float64(time.Second)) // Convert fractional seconds to nanoseconds

	// Create a time.Time object from the Unix timestamp (seconds and nanoseconds)
	// time.Unix always returns a time in UTC if the seconds are from epoch.
	return time.Unix(seconds, nanoseconds), nil
}

// ConvertToUnixTimestamp converts a date string with a specific format
// into a Unix timestamp (in seconds, as a float64 to include fractional seconds).
func ConvertToUnixTimestamp(dateStr string) (float64, error) {
	// Parse the date string using the defined timestamp layout
	t, err := time.Parse(timestampLayout, dateStr)
	if err != nil {
		// Return zero float and error if parsing fails
		return 0.0, fmt.Errorf("failed to parse date string: %w", err)
	}

	// Get the Unix timestamp in nanoseconds and convert to seconds (float64)
	// This is equivalent to Python's dt.timestamp() which returns a float.
	return float64(t.UnixNano()) / float64(time.Second), nil
}

func FromTimeStampToDateStr(timeStamp string) string {
	posicion := 14
	// Verificar que la posición sea válida
	if posicion < 0 || posicion >= len(timeStamp) {
		return timeStamp // No hacer nada si la posición es inválida
	}

	// Dividir la string en dos partes: antes de la posición y después de la posición
	parteAnterior := timeStamp[:posicion]
	partePosterior := timeStamp[posicion:]

	// Concatenar las partes con el punto en medio
	return parteAnterior + "." + partePosterior
}
