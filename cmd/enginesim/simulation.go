package main

import (
	"context"
	"fmt"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"go-playground/internal/justforfun/vehiclesim/differential"
	"go-playground/internal/justforfun/vehiclesim/engine"
	"go-playground/internal/justforfun/vehiclesim/gearbox"
	"go-playground/internal/justforfun/vehiclesim/wheels"
	"log"
	"time"
)

func VehicleSimulation() {
	fmt.Println("Starting vehicle simulation")

	config := ConfigInfluxDB{
		Org:    "docs",
		Bucket: "vehicle-simulation",
	}

	client := NewInfluxDBClient(config)
	defer client.Close()

	writeAPI := client.WriteAPIBlocking(config.Org, config.Bucket)

	// Engine initialization and initial state
	theGearbox := gearbox.NewGearbox()
	theEngine := engine.NewEngine(theGearbox)
	theBasicDifferential := differential.NewBasicDifferential(differential.TypeRDiffRatio)

	wheelManager, err := wheels.NewWheelManager("245/40R19")
	if err != nil {
		panic(fmt.Sprintf("Error initializing wheels: %v", err))
	}

	initializeEngineState(theEngine)
	initializeGearboxState(theGearbox)

	// Simulation Setup
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	// Channel to coordinate the start of goroutines
	ready := make(chan struct{})

	// Accelerator pedal simulation
	go func() {
		<-ready // Wait for the start signal
		simulateEngine(theEngine)
	}()

	// Gear shift simulation
	go func() {
		<-ready // Wait for the start signal
		simulateGearShifts(theEngine, theGearbox)
	}()

	// Signal that the simulation can begin
	fmt.Println("Starting simulation...")
	close(ready)

	// Simulation main loop
	for range ticker.C {
		theEngine.Update(0.1)
		theGearbox.Update(0.1)
		theBasicDifferential.Update(theGearbox.OutputShaft, theGearbox.OutputShaftTorque)

		engineData := theEngine.GetData()
		gearboxData := theGearbox.GetData()
		differentialData := theBasicDifferential.GetData()

		wheelManager.WheelPair.Update(differentialData.WheelSpeedL, differentialData.WheelSpeedR)

		wheelsData := wheelManager.WheelPair.GetData()

		enginePoint := createEnginePoint(engineData)
		gearboxPoint := createGearboxPoint(gearboxData)
		differentialPoint := createDifferentialPoint(differentialData)
		vehicleDynamicPoint := createVehicleDynamicPoint(wheelsData)

		if err := writePoints(writeAPI, enginePoint, gearboxPoint, differentialPoint, vehicleDynamicPoint); err != nil {
			log.Printf("Error writting datas: %v", err)
		}

		printSimulationStatus(engineData, gearboxData, differentialData, wheelsData)

	}
}

func initializeEngineState(motor *engine.Engine) {
	// Initial state of the engine
	motor.SetAcceleratorPos(0.0) // Accelerator depressed

	// Wait for the engine to idle
	fmt.Println("Initializing engine at idle...")
	time.Sleep(2 * time.Second)

	fmt.Printf("Initial state established:\n")
	fmt.Printf("- Engine idling\n")

}

func initializeGearboxState(gearbox *gearbox.Gearbox) {
	// Initial state of the gearbox
	gearbox.SetClutch(0.0) // Clutch pressed
	gearbox.SetGear(0)     // Neutral

	// Wait for the engine to idle
	fmt.Println("Initializing gearbox...")
	time.Sleep(2 * time.Second)

	// Prepare first gear
	gearbox.SetGear(1)     // Engage first gear
	gearbox.SetClutch(1.0) // Release clutch gradually

	fmt.Printf("- First gear engaged\n")
	fmt.Printf("- Clutch ready\n")
}

func simulateEngine(motor *engine.Engine) {
	// Short initial pause for stabilization
	time.Sleep(1 * time.Second)

	for {
		// Smoother gradual acceleration
		for pos := 0.0; pos <= 1.0; pos += 0.05 {
			motor.SetAcceleratorPos(pos)
			time.Sleep(500 * time.Millisecond)
		}
		// Hold full throttle briefly
		time.Sleep(2 * time.Second)

		// Gradual deceleration
		for pos := 1.0; pos >= 0.0; pos -= 0.05 {
			motor.SetAcceleratorPos(pos)
			time.Sleep(500 * time.Millisecond)
		}
		// Idle pause
		time.Sleep(2 * time.Second)
	}
}

func simulateGearShifts(motor *engine.Engine, gearbox *gearbox.Gearbox) {
	// Short initial pause for stabilization
	time.Sleep(2 * time.Second)

	for {
		engineData := motor.GetData()
		gearboxData := gearbox.GetData()

		switch {
		case engineData.RPM > 4000 && gearboxData.CurrentGear < 6:
			performGearShift(motor, gearbox, gearbox.ShiftUp)

		case engineData.RPM < 2000 && gearboxData.CurrentGear > 1:
			performGearShift(motor, gearbox, gearbox.ShiftDown)
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func performGearShift(motor *engine.Engine, gearbox *gearbox.Gearbox, shiftGear func() bool) {
	// Save the current throttle position
	currentAccel := motor.GetData().AcceleratorPosition

	// Gear shift sequence
	motor.SetAcceleratorPos(0.3) // Reduce acceleration
	time.Sleep(100 * time.Millisecond)

	gearbox.SetClutch(0.0) // Press clutch
	time.Sleep(200 * time.Millisecond)

	shiftGear() // Change gear
	time.Sleep(200 * time.Millisecond)

	// Release clutch gradually
	for clutch := 0.0; clutch <= 1.0; clutch += 0.2 {
		gearbox.SetClutch(clutch)
		time.Sleep(50 * time.Millisecond)
	}

	// Restore throttle gradually
	motor.SetAcceleratorPos(currentAccel)
}

func createEnginePoint(engineData engine.Telemetry) *write.Point {
	return write.NewPoint(
		"engine",
		map[string]string{
			"simulation": "engine1",
		},
		map[string]interface{}{
			"rpm":            engineData.RPM,
			"torque":         engineData.Torque,
			"oil_temp":       engineData.OilTemp,
			"accel_position": engineData.AcceleratorPosition,
			"engine_state":   engineData.EngineState,
			"power_kw":       engineData.PowerKW,
			"power_hp":       engineData.PowerHP,
		},
		time.Now(),
	)
}

func createGearboxPoint(gearboxData gearbox.Telemetry) *write.Point {
	return write.NewPoint(
		"gearbox",
		map[string]string{
			"simulation": "gearbox1",
		},
		map[string]interface{}{
			"input_shaft":         gearboxData.InputShaft,
			"output_shaft":        gearboxData.OutputShaft,
			"current_gear":        gearboxData.CurrentGear,
			"clutch_position":     gearboxData.ClutchPosition,
			"input_shaft_torque":  gearboxData.InputShaftTorque,
			"output_shaft_torque": gearboxData.OutputShaftTorque,
			//"wheel_rpm":       gearboxData.WheelRPM,
			//"wheel_torque":    gearboxData.WheelTorque,
		},
		time.Now(),
	)
}

func createDifferentialPoint(differentialData differential.Telemetry) *write.Point {
	return write.NewPoint(
		"differential",
		map[string]string{
			"simulation": "basic_differential",
		},
		map[string]interface{}{
			"wheel_speed_left":  differentialData.WheelSpeedL,
			"wheel_speed_right": differentialData.WheelSpeedR,
		},
		time.Now(),
	)
}

func createVehicleDynamicPoint(wheelData wheels.Telemetry) *write.Point {
	return write.NewPoint(
		"vehicle_dynamic",
		map[string]string{
			"simulation": "vehicle_dynamic",
		},
		map[string]interface{}{
			"vehicle_speed_kmh": wheelData.VehicleSpeed.KMH,
		},
		time.Now(),
	)
}

func writePoints(writeAPI api.WriteAPIBlocking, points ...*write.Point) error {
	for _, point := range points {
		if err := writeAPI.WritePoint(context.Background(), point); err != nil {
			return fmt.Errorf("error writing point: %v", err)
		}
	}
	return nil
}

func printSimulationStatus(engineData engine.Telemetry, gearboxData gearbox.Telemetry, differentialData differential.Telemetry, wheelsData wheels.Telemetry) {
	fmt.Printf(engineData.String())
	fmt.Printf(gearboxData.String())
	fmt.Printf(differentialData.String())
	fmt.Printf(wheelsData.String())
}
