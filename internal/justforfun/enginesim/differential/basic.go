package differential

const TypeRDiffRatio = 3.84

type Differential struct {
	GearRatio   float64 // Differential ratio
	WheelSpeedL float64 // Angular velocity of the left wheel
	WheelSpeedR float64 // Angular velocity of the right wheel
	TorqueL     float64 // Torque sent to the left wheel
	TorqueR     float64 // Torque sent to the right wheel
	SlipRatio   float64 // Slip ratio between the wheels
}

// NewBasicDifferential :
func NewBasicDifferential(gearRatio float64) *Differential {
	return &Differential{
		GearRatio: gearRatio,
	}
}

// Update calculates the speeds and torque for the wheels based on the input
func (d *Differential) Update(inputShaftRPM float64, inputTorque float64, slipRatio float64) {
	// Basic output ratio based on differential ratio
	wheelRPM := inputShaftRPM / d.GearRatio

	// Distribute the torque (basic, no slip)
	d.TorqueL = inputTorque / 2
	d.TorqueR = inputTorque / 2

	// Apply a slip coefficient
	d.WheelSpeedL = wheelRPM * (1.0 - slipRatio/2.0)
	d.WheelSpeedR = wheelRPM * (1.0 + slipRatio/2.0)

	// Update slip ratio
	d.SlipRatio = slipRatio
}
