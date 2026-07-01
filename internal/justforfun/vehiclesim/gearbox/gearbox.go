package gearbox

import "math"

// Gearbox define la interfaz para cualquier sistema de caja de cambios
// Permite que diferentes tipos de cajas (manual, automática, CVT, etc.)
// implementen su propia lógica de transformación RPM/Torque
type Gearbox interface {
	// Update procesa entrada de motor (RPM/Torque) y calcula salida a las ruedas
	// Parámetros:
	//   inputRPM: revoluciones por minuto del eje de entrada (desde motor)
	//   inputTorque: torque en Nm del eje de entrada (desde motor)
	//   deltaTime: tiempo transcurrido en segundos desde la última actualización
	Update(inputRPM float64, inputTorque float64, deltaTime float64)

	// GetData retorna telemetría completa de la caja de cambios
	// Cada implementación retorna su propio tipo concreto (Telemetry, etc.)
	GetData() interface{}

	// GetOutputShaft retorna RPM en eje de salida (hacia diferenciales/ruedas)
	GetOutputShaft() float64

	// GetOutputTorque retorna torque en Nm en eje de salida
	GetOutputTorque() float64

	// SetClutch establece posición del clutch entre 0.0 (disengaged) y 1.0 (engaged)
	SetClutch(position float64)

	// ShiftUp intenta cambio a marcha superior
	// Retorna true si el cambio fue exitoso, false si ya estaba en marcha máxima
	ShiftUp() bool

	// ShiftDown intenta cambio a marcha inferior
	// Retorna true si el cambio fue exitoso, false si ya estaba en marcha mínima
	ShiftDown() bool
}

// GetManualGearboxData es un helper que castea GetData() a Telemetry para ManualGearbox
// Usado internamente para obtener datos tipados
func GetManualGearboxData(g Gearbox) Telemetry {
	if data, ok := g.GetData().(Telemetry); ok {
		return data
	}
	return Telemetry{}
}

// ManualGearbox es la implementación manual de Gearbox
type ManualGearbox struct {
	currentGear    int
	maxGears       int
	gearRatios     []float64
	finalDrive     float64
	ClutchPosition float64 // 0.0 = clutch disengaged, 1.0 = clutch engaged

	InputShaft        float64
	InputShaftTorque  float64
	OutputShaft       float64
	OutputShaftTorque float64

	// Gearbox inertia's (kg·m²)
	inputShaftInertia  float64
	gearInertias       []float64
	outputShaftInertia float64

	// Angular acceleration's (rad/s²)
	inputShaftAcceleration  float64
	outputShaftAcceleration float64
}

// NewManualGearbox crea una nueva instancia de caja de cambios manual
func NewManualGearbox() Gearbox {
	return &ManualGearbox{
		currentGear:    0, // 0 = neutral
		ClutchPosition: 0.0,
		maxGears:       7,
		// Gear ratios
		gearRatios: []float64{
			0.0,   // Neutral
			3.4,   // 1
			2.75,  // 2
			1.767, // 3
			0.925, // 4
			0.705, // 5
			0.755, // 6
			0.635, // 7
		},
		finalDrive: 4.471, // Final differential ratio
		// Initialization of inertias with typical values
		inputShaftInertia: 0.1, // kg·m²
		gearInertias: []float64{
			0.0,   // Neutral
			0.015, // 1
			0.014, // 2
			0.013, // 3
			0.012, // 4
			0.011, // 5
			0.011, // 6
			0.010, // 7
		},
		outputShaftInertia: 0.05, // kg·m²

	}
}

func (g *ManualGearbox) SetClutch(position float64) {
	g.ClutchPosition = math.Max(0, math.Min(1, position))
}

func (g *ManualGearbox) ShiftUp() bool {
	if g.currentGear < g.maxGears {
		g.currentGear++
		return true
	}
	return false
}

func (g *ManualGearbox) ShiftDown() bool {
	if g.currentGear > 0 {
		g.currentGear--
		return true
	}
	return false
}

func (g *ManualGearbox) GetCurrentRatio() float64 {
	return g.gearRatios[g.currentGear] * g.finalDrive
}

// setOutputShaft Calculates output shaft RPM based on input shaft RPM
func (g *ManualGearbox) setOutputShaft(rpm float64) float64 {
	if g.currentGear == 0 {
		return 0
	}
	return rpm / (g.gearRatios[g.currentGear] * g.finalDrive)
}

func (g *ManualGearbox) SetGear(targetGear int) bool {
	if targetGear >= 0 && targetGear <= g.maxGears {
		g.currentGear = targetGear
		return true
	}
	return false
}

// GetOutputShaftTorque Calculate the torque at the wheels
func (g *ManualGearbox) GetOutputShaftTorque(engineTorque float64) float64 {
	efficiency := 0.92 // Transmission efficiency
	return engineTorque * g.GetCurrentRatio() * efficiency * g.ClutchPosition
}

// Function to calculate the inertia of the input shaft
func (g *ManualGearbox) calculateInputShaftInertia() float64 {
	// The inertia of the input shaft is affected by:
	// - Clutch
	// - Engine speed
	clutchEffect := g.ClutchPosition * 0.3 // Clutch influence factor
	baseInertia := g.inputShaftInertia

	return baseInertia * (1 + clutchEffect)
}

// Function to calculate the effective inertia of gears
func (g *ManualGearbox) calculateGearInertia() float64 {
	if g.currentGear == 0 {
		return 0 // In neutral, there is no effective inertia of the gears
	}

	// The effective inertia depends on:
	// - Current gear
	// - The transmission ratio
	// - The basic inertia of the gear
	baseGearInertia := g.gearInertias[g.currentGear]
	gearRatio := g.gearRatios[g.currentGear]

	// The effective inertia increases with the square of the gear ratio.
	return baseGearInertia * math.Pow(gearRatio, 2)
}

// Function to calculate the inertia of the output shaft
func (g *ManualGearbox) calculateOutputShaftInertia() float64 {
	// The inertia of the output shaft is affected by:
	// - OutputShaftSpeed
	// - Current torque
	// - Differential
	baseInertia := g.outputShaftInertia
	loadEffect := math.Abs(g.OutputShaftTorque) * 0.001 // Load factor

	return baseInertia * (1 + loadEffect)
}

// Function to calculate the total inertia of the system
func (g *ManualGearbox) calculateTotalInertia() float64 {
	inputInertia := g.calculateInputShaftInertia()
	gearInertia := g.calculateGearInertia()
	outputInertia := g.calculateOutputShaftInertia()

	return inputInertia + gearInertia + outputInertia
}

// Function to update angular accelerations
func (g *ManualGearbox) updateAngularAccelerations(deltaTime float64) {
	totalInertia := g.calculateTotalInertia()
	if totalInertia > 0 {
		// Calculate angular acceleration based on torque and inertia
		g.inputShaftAcceleration = g.InputShaftTorque / totalInertia
		g.outputShaftAcceleration = g.OutputShaftTorque / totalInertia
	}
}

// Update implementa la interfaz Transmission
// Procesa entrada de motor (inputRPM, inputTorque) y calcula salida a las ruedas
// Parameters:
//
//	inputShaftRPM: revoluciones por minuto del eje de entrada (desde motor)
//	inputShaftTorque: torque en Nm del eje de entrada (desde motor)
//	deltaTime: tiempo transcurrido en segundos
func (g *ManualGearbox) Update(inputShaftRPM float64, inputShaftTorque float64, deltaTime float64) {
	// Store inputs para compatibilidad backwards (debugging)
	g.InputShaft = inputShaftRPM
	g.InputShaftTorque = inputShaftTorque

	// Updates angular accelerations
	g.updateAngularAccelerations(deltaTime)

	if g.ClutchPosition > 0 {
		// Calculate the exit velocity considering the inertias
		baseOutputRPM := g.setOutputShaft(inputShaftRPM)
		inertiaEffect := g.outputShaftAcceleration * deltaTime

		g.OutputShaft = baseOutputRPM + (inertiaEffect * 9.549) // Convert rad/s² to RPM/s
		g.OutputShaftTorque = g.GetOutputShaftTorque(inputShaftTorque)

	} else {
		// Clutch is disengaged, so the output shaft tends to zero
		inertialDeceleration := g.OutputShaft * 0.25 * (1 + g.calculateOutputShaftInertia()*0.1)
		g.OutputShaft = g.OutputShaft - inertialDeceleration
		g.OutputShaftTorque = g.OutputShaftTorque - (g.OutputShaftTorque * 0.25)

	}
}

// GetOutputShaft implementa la interfaz Transmission
// Retorna RPM en eje de salida (hacia diferenciales/ruedas)
func (g *ManualGearbox) GetOutputShaft() float64 {
	return g.OutputShaft
}

// GetOutputTorque implementa la interfaz Transmission
// Retorna torque en Nm en eje de salida (hacia diferenciales/ruedas)
func (g *ManualGearbox) GetOutputTorque() float64 {
	return g.OutputShaftTorque
}

// GetData implementa la interfaz Gearbox
func (g *ManualGearbox) GetData() interface{} {
	return Telemetry{
		ClutchPosition:    g.ClutchPosition,
		InputShaft:        g.InputShaft,
		CurrentGear:       g.currentGear,
		OutputShaft:       g.OutputShaft,
		InputShaftTorque:  g.InputShaftTorque,
		OutputShaftTorque: g.OutputShaftTorque,
	}
}
