package gearbox

import "math"

type Gearbox struct {
	CurrentGear    int
	MaxGears       int
	GearRatios     []float64
	FinalDrive     float64
	ClutchPosition float64 // 0.0 = clutch disengaged, 1.0 = clutch engaged

	InputShaft  float64
	OutputShaft float64

	wheelRPM    float64
	wheelTorque float64
}

func NewGearbox() *Gearbox {
	return &Gearbox{
		CurrentGear:    0, // 0 = neutral
		ClutchPosition: 0.0,
		MaxGears:       6,
		// Gear ratios
		GearRatios: []float64{
			0.0,   // Neutral
			3.827, // 1
			2.359, // 2
			1.656, // 3
			1.221, // 4
			1.000, // 5
			0.831, // 6
		},
		FinalDrive: 3.42, // Final differential ratio
	}
}

func (g *Gearbox) SetClutch(position float64) {
	g.ClutchPosition = math.Max(0, math.Min(1, position))
}

func (g *Gearbox) ShiftUp() bool {
	if g.CurrentGear < g.MaxGears {
		g.CurrentGear++
		return true
	}
	return false
}

func (g *Gearbox) ShiftDown() bool {
	if g.CurrentGear > 0 {
		g.CurrentGear--
		return true
	}
	return false
}

func (g *Gearbox) GetCurrentRatio() float64 {
	return g.GearRatios[g.CurrentGear] * g.FinalDrive
}

// SetOutputShaft Calculates output shaft RPM based on input shaft RPM
func (g *Gearbox) SetOutputShaft(rpm float64) float64 {
	if g.CurrentGear == 0 {
		return 0
	}
	return rpm / (g.GearRatios[g.CurrentGear] * g.FinalDrive)
}

// GetWheelTorque Calculate the torque at the wheels
func (g *Gearbox) GetWheelTorque(engineTorque float64) float64 {
	efficiency := 0.92 // Transmission efficiency
	return engineTorque * g.GetCurrentRatio() * efficiency * g.ClutchPosition
}

func (g *Gearbox) Update(deltaTime float64) {
	// The InputShaft is now updated from the engine

	// Calculate output RPM based on current gear ratio
	g.OutputShaft = g.SetOutputShaft(g.InputShaft)

	// Calculate the torque at the wheels
	if g.ClutchPosition > 0 {
		g.wheelTorque = g.GetWheelTorque(g.InputShaft)
		g.wheelRPM = g.OutputShaft
	} else {
		g.wheelTorque = 0
		g.wheelRPM = 0
	}
}

func (g *Gearbox) GetGearboxData() Telemetry {
	return Telemetry{
		ClutchPosition: g.ClutchPosition,
		InputShaft:     g.InputShaft,
		CurrentGear:    g.CurrentGear,
		OutputShaft:    g.OutputShaft,
		WheelRPM:       g.wheelRPM,
		WheelTorque:    g.wheelTorque,
	}
}

func (g *Gearbox) ShiftGear(gear int) bool {
	if gear >= 0 && gear <= g.MaxGears {
		g.CurrentGear = gear
		return true
	}
	return false
}
