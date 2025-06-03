package gearbox

import "math"

type Gearbox struct {
	CurrentGear    int
	MaxGears       int
	GearRatios     []float64
	FinalDrive     float64
	ClutchPosition float64 // 0.0 = clutch disengaged, 1.0 = clutch engaged
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

// GetWheelRPM Calculates wheel RPM based on engine RPM
func (g *Gearbox) GetWheelRPM(engineRPM float64) float64 {
	if g.CurrentGear == 0 {
		return 0
	}
	return engineRPM / (g.GearRatios[g.CurrentGear] * g.FinalDrive)
}

// GetWheelTorque Calculate the torque at the wheels
func (g *Gearbox) GetWheelTorque(engineTorque float64) float64 {
	efficiency := 0.92 // Transmission efficiency
	return engineTorque * g.GetCurrentRatio() * efficiency * g.ClutchPosition
}
