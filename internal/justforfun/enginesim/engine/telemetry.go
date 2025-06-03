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

// GetPowerHP returns the power in horsepower
func (d Telemetry) GetPowerHP() float64 {
	return d.PowerKW * 1.34102
}

// GetEfficiency calculates an approximate efficiency
func (d Telemetry) GetEfficiency() float64 {
	maxTheoreticalPower := 150.0 // kW, adjust to specifications
	return (d.PowerKW / maxTheoreticalPower) * 100
}

// String implements the String interface for human-readable formatting
func (d Telemetry) String() string {
	return fmt.Sprintf(
		"Engine [RPM: %.0f, torque: %.1f Nm, Temp: %.1f°C, Pot: %.1f kW, Pot: %.1f HP, State: %s]",
		d.RPM, d.Torque, d.OilTemp, d.PowerKW, d.PowerHP, d.EngineState,
	)
}
