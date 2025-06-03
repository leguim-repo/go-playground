package main

import (
	"context"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"go-playground/internal/justforfun/enginesim/engine"
	"log"
	"os"
	"time"
)

type ConfigInfluxDB struct {
	URL    string
	Token  string
	Org    string
	Bucket string
}

// NewInfluxDBClient create a new InfluxDB client
func NewInfluxDBClient(config ConfigInfluxDB) influxdb2.Client {
	if config.Token == "" {
		config.Token = os.Getenv("INFLUXDB_TOKEN")
	}

	if config.URL == "" {
		config.URL = "http://localhost:8086"
	}

	return influxdb2.NewClient(config.URL, config.Token)
}

func mainVisual() {
	config := ConfigInfluxDB{
		Org:    "docs",
		Bucket: "engine-torque-curve",
	}

	client := NewInfluxDBClient(config)
	defer client.Close()

	writeAPI := client.WriteAPIBlocking(config.Org, config.Bucket)

	// Generate torque curves for different throttle positions
	acceleratorPositions := []float64{0.25, 0.5, 0.75, 1.0}

	motor := engine.NewEngine()

	for _, position := range acceleratorPositions {
		motor.SetAccelerator(position)

		// Generate points throughout the RPM range
		for rpm := 800.0; rpm <= motor.MaxRPM; rpm += 100 {
			motor.Rpm = rpm
			motor.UpdateTorque()

			engineData := motor.GetData()

			// Create a point for InfluxDB
			point := write.NewPoint(
				"torque_curve",
				map[string]string{
					"simulation": "engine1",
					"accelPos":   fmt.Sprintf("%.2f", position),
				},
				map[string]interface{}{
					"rpm":      engineData.RPM,
					"torque":   engineData.Torque,
					"power_kw": engineData.PowerKW, // Power in kW
				},
				time.Now(),
			)

			if err := writeAPI.WritePoint(context.Background(), point); err != nil {
				log.Printf("Error writing data: %v", err)
			}
		}
	}
}

func mainSim() {
	config := ConfigInfluxDB{
		Org:    "docs",
		Bucket: "engine-simulation",
	}

	client := NewInfluxDBClient(config)
	defer client.Close()

	writeAPI := client.WriteAPIBlocking(config.Org, config.Bucket)

	motor := engine.NewEngine()

	// Driving simulation
	ticker := time.NewTicker(100 * time.Millisecond) // Update every 100ms
	defer ticker.Stop()

	// Simulate accelerator pedal movement
	go func() {
		for {
			for pos := 0.0; pos <= 1.0; pos += 0.1 {
				motor.SetAccelerator(pos)
				time.Sleep(1 * time.Second)
			}
			for pos := 1.0; pos >= 0.0; pos -= 0.1 {
				motor.SetAccelerator(pos)
				time.Sleep(1 * time.Second)
			}
		}
	}()

	// Main loop simulation
	for range ticker.C {
		motor.Update(0.1) // deltaTime = 0.1 seconds
		engineData := motor.GetData()
		fmt.Printf("engineData: %+v\n", engineData)

		point := write.NewPoint(
			"engine_data",
			map[string]string{
				"simulation": "engine1",
			},
			map[string]interface{}{
				"rpm":      engineData.RPM,
				"torque":   engineData.Torque,
				"oilTemp":  engineData.OilTemp,
				"accelPos": engineData.AcceleratorPosition,
			},
			time.Now(),
		)

		if err := writeAPI.WritePoint(context.Background(), point); err != nil {
			log.Printf("Error escribiendo datos: %v", err)
		}
	}
}

func main() {
	mainVisual()
	mainSim()
}
