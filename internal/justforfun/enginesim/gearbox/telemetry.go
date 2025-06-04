package gearbox

import "fmt"

type Telemetry struct {
	CurrentGear       int
	ClutchPosition    float64
	InputShaft        float64
	InputShaftTorque  float64
	OutputShaft       float64
	OutputShaftTorque float64
}

// String implements the String interface for human-readable formatting
func (d Telemetry) String() string {
	return fmt.Sprintf(
		"Gearbox [Gear=%d, Clutch=%.1f %s, InputShaft: %.0f rpm, OutputShaft: %.0f rpm, InputShaftTorque=%.1f Nm, OutputShaftTorque=%.1f Nm]\n",
		d.CurrentGear,
		d.getClutchPositionPercentile(),
		" %%", // Separate the percentage to avoid errors
		d.InputShaft,
		d.OutputShaft,
		d.InputShaftTorque,
		d.OutputShaftTorque,
	)
}

func (d Telemetry) getClutchPositionPercentile() float64 {
	return d.ClutchPosition * 100
}
