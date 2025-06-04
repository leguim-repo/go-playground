package wheels

import "fmt"

// Telemetry provides complete wheel telemetry data
type Telemetry struct {
	WheelSpeedL  float64
	WheelSpeedR  float64
	VehicleSpeed Speed
	TireInfo     TireInfo
}

func (d Telemetry) String() string {
	return fmt.Sprintf("Wheels [WheelSpeedL: %.2f RPM, WheelSpeedR: %.2f RPM, VehicleSpeed: %.2f KMH]\n",
		d.WheelSpeedL,
		d.WheelSpeedR,
		d.VehicleSpeed.KMH)

}
