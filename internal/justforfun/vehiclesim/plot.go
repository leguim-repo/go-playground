package vehiclesim

import (
	"context"
	"fmt"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"go-playground/internal/justforfun/vehiclesim/engine"
	"go-playground/internal/justforfun/vehiclesim/gearbox"
	"go-playground/internal/justforfun/vehiclesim/influx"
	"log"
	"time"
)

func PlotEngineTorqueCurve() {
	fmt.Println("Plotting engine torque curve")
	config := influx.ConfigInfluxDB{
		Org:    "docs",
		Bucket: "engine-torque-curve",
	}

	client := influx.NewInfluxDBClient(config)
	defer client.Close()

	writeAPI := client.WriteAPIBlocking(config.Org, config.Bucket)

	// Generate torque curves for different throttle positions
	acceleratorPositions := []float64{0.25, 0.5, 0.75, 1.0}

	motor := engine.NewEngine(gearbox.NewGearbox())

	for _, position := range acceleratorPositions {
		motor.SetAcceleratorPos(position)

		// Generate points throughout the RPM range
		for rpm := 800.0; rpm <= motor.MaxRPM; rpm += 100 {
			motor.Rpm = rpm
			motor.UpdateTorque()

			engineData := motor.GetData()

			// Create a point for InfluxDB
			point := write.NewPoint(
				"torque_curve",
				map[string]string{
					"simulation":     "engine1",
					"accel_position": fmt.Sprintf("%.2f", position),
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
