package wheels

import "math"

// TireInfo contains detailed tire information
type TireInfo struct {
	WidthMM          float64
	AspectRatio      float64
	WheelDiameterIn  float64
	SideWallHeightMM float64
	TotalRadiusM     float64
	CircumferenceM   float64
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
