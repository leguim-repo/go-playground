package gearbox

import "math"

type Gearbox struct {
	currentGear    int
	maxGears       int
	gearRatios     []float64
	finalDrive     float64
	ClutchPosition float64 // 0.0 = clutch disengaged, 1.0 = clutch engaged

	InputShaft        float64
	InputShaftTorque  float64
	outputShaft       float64
	outputShaftTorque float64

	wheelRPM float64
}

func NewGearbox() *Gearbox {
	return &Gearbox{
		currentGear:    0, // 0 = neutral
		ClutchPosition: 0.0,
		maxGears:       6,
		// Gear ratios
		gearRatios: []float64{
			0.0,   // Neutral
			3.827, // 1
			2.359, // 2
			1.656, // 3
			1.221, // 4
			1.000, // 5
			0.831, // 6
		},
		finalDrive: 3.42, // Final differential ratio
	}
}

func (g *Gearbox) SetClutch(position float64) {
	g.ClutchPosition = math.Max(0, math.Min(1, position))
}

func (g *Gearbox) ShiftUp() bool {
	if g.currentGear < g.maxGears {
		g.currentGear++
		return true
	}
	return false
}

func (g *Gearbox) ShiftDown() bool {
	if g.currentGear > 0 {
		g.currentGear--
		return true
	}
	return false
}

func (g *Gearbox) GetCurrentRatio() float64 {
	return g.gearRatios[g.currentGear] * g.finalDrive
}

// setOutputShaft Calculates output shaft RPM based on input shaft RPM
func (g *Gearbox) setOutputShaft(rpm float64) float64 {
	if g.currentGear == 0 {
		return 0
	}
	return rpm / (g.gearRatios[g.currentGear] * g.finalDrive)
}

func (g *Gearbox) Update(deltaTime float64) {
	// The InputShaft and InputShaftTorque is now updated from the engine. Check TODO of engine

	// Calculate output RPM based on current gear ratio
	g.outputShaft = g.setOutputShaft(g.InputShaft)

	// Calculate the torque at the wheels
	if g.ClutchPosition > 0 {
		g.outputShaftTorque = g.GetOutputShaftTorque(g.InputShaftTorque)
		g.wheelRPM = g.outputShaft
	} else {
		g.wheelRPM = 0
	}
}

func (g *Gearbox) SetGear(targetGear int) bool {
	if targetGear >= 0 && targetGear <= g.maxGears {
		g.currentGear = targetGear
		return true
	}
	return false
}

func (g *Gearbox) GetGearboxData() Telemetry {
	return Telemetry{
		ClutchPosition:    g.ClutchPosition,
		InputShaft:        g.InputShaft,
		CurrentGear:       g.currentGear,
		OutputShaft:       g.outputShaft,
		OutputShaftTorque: g.outputShaftTorque,
	}
}

// GetOutputShaftTorque Calculate the torque at the wheels
func (g *Gearbox) GetOutputShaftTorque(engineTorque float64) float64 {
	efficiency := 0.92 // Transmission efficiency
	return engineTorque * g.GetCurrentRatio() * efficiency * g.ClutchPosition
}
