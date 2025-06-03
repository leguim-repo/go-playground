package gearbox

type Telemetry struct {
	ClutchPosition float64
	InputShaft     float64
	CurrentGear    int
	OutputShaft    float64
	WheelRPM       float64
	WheelTorque    float64
}
