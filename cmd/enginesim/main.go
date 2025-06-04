package main

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"os"
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

func main() {
	PlotEngineTorqueCurve()
	VehicleSimulation()

}
