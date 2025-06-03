package differential

import "fmt"

type Telemetry struct {
	WheelSpeedL float64
	WheelSpeedR float64
}

func (d *Differential) GetTelemetry() Telemetry {
	return Telemetry{
		WheelSpeedL: d.WheelSpeedL,
		WheelSpeedR: d.WheelSpeedR,
	}
}

func (d *Differential) String() string {
	return fmt.Sprintf("Differential [L: %.0f RPM, R: %.0f RPM],",
		d.WheelSpeedL,
		d.WheelSpeedR)
}
