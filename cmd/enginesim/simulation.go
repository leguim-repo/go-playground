package main

import (
	"context"
	"fmt"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"go-playground/internal/justforfun/enginesim/engine"
	"go-playground/internal/justforfun/enginesim/gearbox"
	"log"
	"time"
)

func EngineSimulation() {
	fmt.Println("Starting engine simulation")

	config := ConfigInfluxDB{
		Org:    "docs",
		Bucket: "engine-simulation",
	}

	client := NewInfluxDBClient(config)
	defer client.Close()

	writeAPI := client.WriteAPIBlocking(config.Org, config.Bucket)

	// Engine initialization and initial state
	motor := engine.NewEngine()
	initializeEngineState(motor)

	// Simulation Setup
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	// Channel to coordinate the start of goroutines
	ready := make(chan struct{})

	// Accelerator pedal simulation
	go func() {
		<-ready // Wait for the start signal
		simulateAccelerator(motor)
	}()

	// Gear shift simulation
	go func() {
		<-ready // Wait for the start signal
		simulateGearShifts(motor)
	}()

	// Signal that the simulation can begin
	fmt.Println("Starting simulation...")
	close(ready)

	// Simulation main loop
	for range ticker.C {
		motor.Update(0.1)
		engineData := motor.GetData()
		gearboxData := motor.GetGearboxData()

		enginePoint := createEnginePoint(engineData)
		gearboxPoint := createGearboxPoint(gearboxData)

		if err := writePoints(writeAPI, enginePoint, gearboxPoint); err != nil {
			log.Printf("Error writting datas: %v", err)
		}

		printSimulationStatus(engineData, gearboxData)
	}
}

func initializeEngineState(motor *engine.Engine) {
	// Initial state of the engine
	motor.SetAccelerator(0.0) // Accelerator depressed

	// Initial state of the gearbox
	motor.SetClutch(0.0) // Clutch pressed
	motor.ShiftGear(0)   // Neutral

	// Wait for the engine to idle
	fmt.Println("Initializing engine at idle...")
	time.Sleep(2 * time.Second)

	// Prepare first gear
	motor.ShiftGear(1)   // First gear
	motor.SetClutch(1.0) // Release clutch gradually

	fmt.Printf("Initial state established:\n")
	fmt.Printf("- Engine idling\n")
	fmt.Printf("- First gear engaged\n")
	fmt.Printf("- Clutch ready\n")
}

func simulateAccelerator(motor *engine.Engine) {
	// Short initial pause for stabilization
	time.Sleep(1 * time.Second)

	for {
		// Smoother gradual acceleration
		for pos := 0.0; pos <= 1.0; pos += 0.05 {
			motor.SetAccelerator(pos)
			time.Sleep(500 * time.Millisecond)
		}
		// Hold full throttle briefly
		time.Sleep(2 * time.Second)

		// Gradual deceleration
		for pos := 1.0; pos >= 0.0; pos -= 0.05 {
			motor.SetAccelerator(pos)
			time.Sleep(500 * time.Millisecond)
		}
		// Idle pause
		time.Sleep(2 * time.Second)
	}
}

func simulateGearShifts(motor *engine.Engine) {
	// Short initial pause for stabilization
	time.Sleep(2 * time.Second)

	for {
		engineData := motor.GetData()
		gearData := motor.GetGearboxData()

		switch {
		case engineData.RPM > 4000 && gearData.CurrentGear < 6:
			performGearShift(motor, gearData.CurrentGear+1)

		case engineData.RPM < 2000 && gearData.CurrentGear > 1:
			performGearShift(motor, gearData.CurrentGear-1)
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func performGearShift(motor *engine.Engine, targetGear int) {
	// Save the current throttle position
	currentAccel := motor.GetData().AcceleratorPosition

	// Gear shift sequence
	motor.SetAccelerator(0.3) // Reduce acceleration
	time.Sleep(100 * time.Millisecond)

	motor.SetClutch(0.0) // Press clutch
	time.Sleep(200 * time.Millisecond)

	motor.ShiftGear(targetGear) // Change gear
	time.Sleep(200 * time.Millisecond)

	// Release clutch gradually
	for clutch := 0.0; clutch <= 1.0; clutch += 0.2 {
		motor.SetClutch(clutch)
		time.Sleep(50 * time.Millisecond)
	}

	// Restore throttle gradually
	motor.SetAccelerator(currentAccel)
}

func createEnginePoint(engineData engine.Telemetry) *write.Point {
	return write.NewPoint(
		"engine_data",
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
		"gearbox_data",
		map[string]string{
			"simulation": "engine1",
		},
		map[string]interface{}{
			"current_gear":    gearboxData.CurrentGear,
			"clutch_position": gearboxData.ClutchPosition,
			"wheel_rpm":       gearboxData.WheelRPM,
			"wheel_torque":    gearboxData.WheelTorque,
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

func printSimulationStatus(engineData engine.Telemetry, gearboxData gearbox.Telemetry) {
	fmt.Printf("Engine: RPM=%.0f, Torque=%.1f Nm, Gear=%d, Clutch=%.1f%%\n",
		engineData.RPM,
		engineData.Torque,
		gearboxData.CurrentGear,
		gearboxData.ClutchPosition*100)
	fmt.Printf("Wheels: RPM=%.0f, Torque=%.1f Nm\n",
		gearboxData.WheelRPM,
		gearboxData.WheelTorque)
}
