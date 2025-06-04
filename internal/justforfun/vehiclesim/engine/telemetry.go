package engine

import "fmt"

type Telemetry struct {
	RPM                 float64
	Torque              float64
	OilTemp             float64
	AcceleratorPosition float64
	PowerKW             float64
	PowerHP             float64
	EngineState         string
}

// String implements the String interface for human-readable formatting
func (d Telemetry) String() string {
	return fmt.Sprintf(
		"Engine [Speed: %.0f, AcelPos: %.1f %s, torque: %.1f Nm, OilTemp: %.1f°C, Power: %.1f kW, Power: %.1f HP, State: %s]\n",
		d.RPM,
		d.getAcceleratorPositionPercentile(),
		" %%",
		d.Torque,
		d.OilTemp,
		d.PowerKW,
		d.PowerHP,
		d.EngineState,
	)
}

func (d Telemetry) getAcceleratorPositionPercentile() float64 {
	return d.AcceleratorPosition * 100
}
