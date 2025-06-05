package differential

const TypeRDiffRatio = 3.84

type Differential struct {
	gearRatio   float64 // Differential ratio
	wheelSpeedL float64 // Angular velocity of the left wheel
	wheelSpeedR float64 // Angular velocity of the right wheel
	torqueL     float64 // Torque sent to the left wheel
	torqueR     float64 // Torque sent to the right wheel
	slipRatio   float64 // Slip ratio between the wheels
}

// NewBasicDifferential :
func NewBasicDifferential(gearRatio float64) *Differential {
	return &Differential{
		gearRatio: gearRatio,
	}
}

// Update calculates the speeds and torque for the wheels based on the input
func (d *Differential) Update(inputShaftRPM float64, inputTorque float64, slipRatio float64) {
	// slipRatio is an external input and should be calculated in based on steering angle, terrain, wheel grip, etc.

	// Basic output ratio based on a differential ratio
	wheelRPM := inputShaftRPM / d.gearRatio

	// Distribute the torque (basic, no slip)
	d.torqueL = inputTorque / 2
	d.torqueR = inputTorque / 2

	// Apply a slip coefficient
	d.wheelSpeedL = wheelRPM * (1.0 - slipRatio/2.0)
	d.wheelSpeedR = wheelRPM * (1.0 + slipRatio/2.0)

	// Update slip ratio
	d.slipRatio = slipRatio
}

func (d *Differential) GetData() Telemetry {
	return Telemetry{
		WheelSpeedL: d.wheelSpeedL,
		WheelSpeedR: d.wheelSpeedR,
	}
}
