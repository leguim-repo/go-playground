package engine

import (
	"go-playground/internal/justforfun/enginesim/gearbox"
	"math"
	"math/rand"
)

type Engine struct {
	// Engine state
	Rpm             float64
	torque          float64
	oilTemp         float64
	acceleratorPos  float64 // 0.0 to 1.0 (0% to 100%)
	fuelConsumption float64
	oilPressure     float64
	waterTemp       float64

	// Engine limits
	MaxRPM    float64
	maxTorque float64
	maxTemp   float64
	minTemp   float64

	// Engine dynamics
	inertia float64 // How quickly the Engine responds

	// torque curve parameters
	rpmMaxTorque float64 // RPM where maximum torque is reached
	rpmMaxPower  float64 // RPM where maximum power is reached

	Gearbox *gearbox.Gearbox
}

func NewEngine(theGearbox *gearbox.Gearbox) *Engine {

	return &Engine{
		Gearbox:        theGearbox,
		Rpm:            800, // Low Idle
		torque:         0,
		oilTemp:        80, // Initial oil temperature
		acceleratorPos: 0,
		MaxRPM:         8500, // Max RPM
		maxTorque:      450,  // Nm
		maxTemp:        120,  // Max oil temperature
		minTemp:        70,   // Min oil temperature in normal conditions
		inertia:        0.3,  // Inertia Engine factor (0-1)
		rpmMaxTorque:   3500, // Typical RPM for maximum torque
		rpmMaxPower:    5500, // Typical RPM for maximum power
	}
}

func randomInRange(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

func (m *Engine) SetAcceleratorPos(position float64) {
	// Ensure a valid position between 0 and 1
	m.acceleratorPos = math.Max(0, math.Min(1, position))
}

func (m *Engine) Update(deltaTime float64) {

	// Update RPM considering the clutch
	clutchSlip := 1.0 - m.Gearbox.ClutchPosition

	// If the clutch is pressed (clutch disengaged), the engine spins more freely.
	rpmDrop := m.Rpm * clutchSlip * 0.1

	m.updateRPM(deltaTime)
	m.Rpm -= rpmDrop

	m.UpdateTorque()
	m.updateOilTemp(deltaTime)

	// TODO: Improve coupling between engine and gearbox. This way is the best way? Maybe a better way is to have a channels between the engine goroutine and the gearbox goroutine.
	// Coupling engine with the gearbox. The rpm engine is the gearbox input shaft
	m.Gearbox.InputShaft = m.Rpm
	// Update gearbox
	m.Gearbox.Update(deltaTime)
}

func (m *Engine) updateRPM(deltaTime float64) {
	// Calculate target RPM based on throttle position
	rpmTarget := m.acceleratorPos*(m.MaxRPM-800) + 800

	// Add random variation to simulate fluctuations
	noise := randomInRange(-50, 50)

	// Interpolate smoothly towards the target using inertia
	m.Rpm = m.Rpm + (rpmTarget-m.Rpm)*m.inertia*deltaTime + noise

	// Limit RPM. To cut!!
	m.Rpm = math.Max(800, math.Min(m.MaxRPM, m.Rpm))
}

func (m *Engine) realisticTorqueCurve(rpm float64) float64 {
	// Normalize RPM to the 0-1 range
	rpmNorm := rpm / m.MaxRPM

	// Parameters to adjust the shape of the curve
	torqueMaxRPM := m.rpmMaxTorque / m.MaxRPM
	powerMaxRPM := m.rpmMaxPower / m.MaxRPM

	// Create a curve that:
	// - Starts low at idle
	// - Rises rapidly
	// - Peaks at rpmMaxTorque
	// - Gradually drops to rpmMaxPower
	// - Drops more rapidly thereafter

	// Principal component of the curve
	baseCurve := math.Exp(-math.Pow((rpmNorm-torqueMaxRPM)*2.5, 2))

	// Add a drop to high RPM
	highDrop := math.Exp(-math.Pow((rpmNorm-powerMaxRPM)*1.5, 2))

	// Reduce torque at idle
	idleFactor := 1 - math.Exp(-5*rpmNorm)

	// Combine all the factors
	torqueFactor := baseCurve * idleFactor * (0.7 + 0.3*highDrop)

	// Multiply by the maximum torque and the throttle position
	return torqueFactor * m.maxTorque * m.acceleratorPos
}

func (m *Engine) UpdateTorque() {
	m.torque = m.realisticTorqueCurve(m.Rpm)

	// Add a small random variation (1-2% of current torque)
	smallRandomTorqueVariation := m.torque * randomInRange(-0.02, 0.02)
	m.torque += smallRandomTorqueVariation

	// Make sure it is not negative
	m.torque = math.Max(0, m.torque)

}

func (m *Engine) updateOilTemp(deltaTime float64) {
	// Temperature increases with RPM and load
	tempTarget := m.minTemp +
		(m.maxTemp-m.minTemp)*(0.3*m.Rpm/m.MaxRPM+0.7*m.acceleratorPos)

	// Add random variation
	noise := randomInRange(-0.5, 0.5)

	// Gradual change in temperature
	m.oilTemp += (tempTarget-m.oilTemp)*0.1*deltaTime + noise

	// Limit temperature
	m.oilTemp = math.Max(m.minTemp, math.Min(m.maxTemp, m.oilTemp))
}

// GetData Function to collect data from the Engine
func (m *Engine) GetData() Telemetry {
	powerKW := (m.torque * m.Rpm) / 9549.297
	powerHP := powerKW * 1.341

	return Telemetry{
		RPM:                 m.Rpm,
		Torque:              m.torque,
		OilTemp:             m.oilTemp,
		AcceleratorPosition: m.acceleratorPos,
		PowerKW:             powerKW,
		PowerHP:             powerHP,
		EngineState:         m.getState(),
	}

}

// randomEngineEvents Function to simulate random engine events
func (m *Engine) randomEngineEvents() string {
	// 0.1% chance
	if rand.Float64() < 0.001 {
		events := []string{
			"temperature_fluctuation",
			"oil_pressure_drop",
			"abnormal_vibration",
		}
		return events[rand.Intn(len(events))]
	}
	return "normal"
}

func (m *Engine) getState() string {
	switch {
	case m.Rpm < 850:
		return "low_idle"
	case m.Rpm >= m.MaxRPM*0.95:
		return "rpm_limit"
	case m.oilTemp >= m.maxTemp*0.9:
		return "oilTemp_high"
	default:
		return "normal"
	}
}
