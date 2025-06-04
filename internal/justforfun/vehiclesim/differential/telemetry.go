package differential

import "fmt"

type Telemetry struct {
	WheelSpeedL float64
	WheelSpeedR float64
}

func (d Telemetry) String() string {
	return fmt.Sprintf("Differential [WheelSpeedL: %.0f RPM, WheelSpeedR: %.0f RPM]\n",
		d.WheelSpeedL,
		d.WheelSpeedR)
}
