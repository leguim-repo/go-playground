package wheels

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

// Update updates wheel speeds
func (wp *WheelPair) Update(leftRPM, rightRPM float64) {
	wp.Left.SetSpeedRPM(leftRPM)
	wp.Right.SetSpeedRPM(rightRPM)
}

// GetData returns complete telemetry data
func (wp *WheelPair) GetData() Telemetry {
	return Telemetry{
		WheelSpeedL:  wp.Left.GetSpeedRPM(),
		WheelSpeedR:  wp.Right.GetSpeedRPM(),
		VehicleSpeed: wp.GetVehicleSpeed(),
		TireInfo:     wp.Left.GetTireInfo(), // Assuming same size on both wheels
	}
}
