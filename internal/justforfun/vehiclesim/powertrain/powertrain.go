package powertrain

import (
	"go-playground/internal/justforfun/vehiclesim/engine"
	"go-playground/internal/justforfun/vehiclesim/gearbox"
)

// PowertrainController orquesta la integración del Motor con la Caja de Cambios de forma integrada
// Encapsula la lógica de acoplamiento entre componentes en una sola abstracción
type PowertrainController struct {
	engine  *engine.Engine
	gearbox gearbox.Gearbox
}

// NewPowertrainController crea una nueva instancia del controlador de transmisión
// Parameters:
//
//	eng: instancia del motor
//	gb: implementación de la interfaz Gearbox (ej: ManualGearbox)
func NewPowertrainController(
	eng *engine.Engine,
	gb gearbox.Gearbox,
) *PowertrainController {
	return &PowertrainController{
		engine:  eng,
		gearbox: gb,
	}
}

// Update orquesta la actualización completa: motor + caja de cambios
// Mantiene los componentes sincronizados propagando datos correctamente
// Parameters:
//
//	clutchPosition: posición del clutch (0.0-1.0)
//	deltaTime: tiempo transcurrido en segundos
func (pc *PowertrainController) Update(clutchPosition float64, deltaTime float64) {
	// Actualizar motor con posición del clutch
	pc.engine.Update(clutchPosition, deltaTime)

	// Propagar datos del motor a la caja de cambios
	pc.gearbox.Update(pc.engine.GetRPM(), pc.engine.GetTorque(), deltaTime)
}

// GetEngineData retorna datos telemétricos del motor
func (pc *PowertrainController) GetEngineData() engine.Telemetry {
	return pc.engine.GetData()
}

// GetGearboxData retorna datos telemétricos de la caja de cambios
// Castea a Telemetry si la implementación es ManualGearbox
func (pc *PowertrainController) GetGearboxData() interface{} {
	return pc.gearbox.GetData()
}

// GetManualGearboxDataTyped es un helper que retorna datos tipados de la caja
func (pc *PowertrainController) GetManualGearboxDataTyped() gearbox.Telemetry {
	return gearbox.GetManualGearboxData(pc.gearbox)
}
