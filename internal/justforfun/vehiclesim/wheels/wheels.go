package wheels

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
)

const (
	// Unit conversions
	RPMToRadPerSec = 2 * math.Pi / 60
	MSToKMH        = 3.6
	MSToMPH        = 2.237
	InchToMM       = 25.4
	MMToM          = 0.001
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

// Wheel represents a wheel with its tire
type Wheel struct {
	tireSize *TireSize
	speedRPM float64
}

// NewWheel creates a new Wheel instance with specific tire size
func NewWheel(tireSpec string) (*Wheel, error) {
	tireSize, err := ParseTireSize(tireSpec)
	if err != nil {
		return nil, err
	}

	return &Wheel{
		tireSize: tireSize,
		speedRPM: 0,
	}, nil
}

// TireInfo contains detailed tire information
type TireInfo struct {
	WidthMM          float64
	AspectRatio      float64
	WheelDiameterIn  float64
	SideWallHeightMM float64
	TotalRadiusM     float64
	CircumferenceM   float64
}

// GetTireInfo returns detailed tire information
func (w *Wheel) GetTireInfo() TireInfo {
	return TireInfo{
		WidthMM:          w.tireSize.Width,
		AspectRatio:      w.tireSize.AspectRatio,
		WheelDiameterIn:  w.tireSize.WheelDiameter,
		SideWallHeightMM: w.tireSize.SideWallHeight,
		TotalRadiusM:     w.tireSize.TotalRadius,
		CircumferenceM:   2 * math.Pi * w.tireSize.TotalRadius,
	}
}

// SetSpeedRPM sets wheel rotation speed in RPM
func (w *Wheel) SetSpeedRPM(rpm float64) {
	w.speedRPM = rpm
}

// GetSpeedRPM gets current wheel rotation speed in RPM
func (w *Wheel) GetSpeedRPM() float64 {
	return w.speedRPM
}

// GetLinearSpeedMS calculates linear speed in meters per second
func (w *Wheel) GetLinearSpeedMS() float64 {
	angularVelocity := w.speedRPM * RPMToRadPerSec
	return angularVelocity * w.tireSize.TotalRadius
}

// WheelPair represents a pair of wheels (left and right)
type WheelPair struct {
	Left  *Wheel
	Right *Wheel
}

// NewWheelPair creates a new pair of wheels with specified tire size
func NewWheelPair(tireSpec string) (*WheelPair, error) {
	left, err := NewWheel(tireSpec)
	if err != nil {
		return nil, err
	}

	right, err := NewWheel(tireSpec)
	if err != nil {
		return nil, err
	}

	return &WheelPair{
		Left:  left,
		Right: right,
	}, nil
}

// Speed represents velocity in different units
type Speed struct {
	MS       float64 // Meters per second
	KMH      float64 // Kilometers per hour
	MPH      float64 // Miles per hour
	WheelRPM float64 // Wheel RPM
}

// Telemetry provides complete wheel telemetry data
type Telemetry struct {
	WheelSpeedL  float64
	WheelSpeedR  float64
	VehicleSpeed Speed
	TireInfo     TireInfo
}

// GetVehicleSpeed calculates vehicle speed based on wheel speeds
func (wp *WheelPair) GetVehicleSpeed() Speed {
	leftSpeed := wp.Left.GetLinearSpeedMS()
	rightSpeed := wp.Right.GetLinearSpeedMS()
	avgSpeedMS := (leftSpeed + rightSpeed) / 2

	return Speed{
		MS:       avgSpeedMS,
		KMH:      avgSpeedMS * MSToKMH,
		MPH:      avgSpeedMS * MSToMPH,
		WheelRPM: (wp.Left.GetSpeedRPM() + wp.Right.GetSpeedRPM()) / 2,
	}
}

// GetTelemetry returns complete telemetry data
func (wp *WheelPair) GetTelemetry() Telemetry {
	return Telemetry{
		WheelSpeedL:  wp.Left.GetSpeedRPM(),
		WheelSpeedR:  wp.Right.GetSpeedRPM(),
		VehicleSpeed: wp.GetVehicleSpeed(),
		TireInfo:     wp.Left.GetTireInfo(), // Assuming same size on both wheels
	}
}

// Update updates wheel speeds
func (wp *WheelPair) Update(leftRPM, rightRPM float64) {
	wp.Left.SetSpeedRPM(leftRPM)
	wp.Right.SetSpeedRPM(rightRPM)
}
