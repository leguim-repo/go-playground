package wheels

import (
	"fmt"
	"regexp"
	"strconv"
)

// TireSize represents tire dimensions
type TireSize struct {
	Width          float64 // Width in millimeters
	AspectRatio    float64 // Aspect ratio (height/width)
	WheelDiameter  float64 // Rim diameter in inches
	SideWallHeight float64 // Sidewall height in millimeters
	TotalRadius    float64 // Total radius in meters
}

// ParseTireSize parses a tire specification (e.g., "245/40R19")
// Format: [Width]/[AspectRatio]R[WheelDiameter]
func ParseTireSize(spec string) (*TireSize, error) {
	pattern := regexp.MustCompile(`^(\d+)/(\d+)R(\d+)$`)
	matches := pattern.FindStringSubmatch(spec)

	if matches == nil {
		return nil, fmt.Errorf("invalid tire format: %s", spec)
	}

	// Convert string values to numbers
	width, _ := strconv.ParseFloat(matches[1], 64)
	aspectRatio, _ := strconv.ParseFloat(matches[2], 64)
	diameter, _ := strconv.ParseFloat(matches[3], 64)

	// Calculate dimensions
	sideWallHeight := width * (aspectRatio / 100)
	wheelDiameterMM := diameter * InchToMM
	totalDiameterMM := wheelDiameterMM + (2 * sideWallHeight)

	return &TireSize{
		Width:          width,
		AspectRatio:    aspectRatio,
		WheelDiameter:  diameter,
		SideWallHeight: sideWallHeight,
		TotalRadius:    (totalDiameterMM / 2) * MMToM, // Convert to meters
	}, nil
}
