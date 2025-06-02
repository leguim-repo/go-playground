package influxdb

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

func randomInRangeInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func randomInRangeFloat(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

func PlaygroundInfluxDbWriter() {
	token := os.Getenv("INFLUXDB_TOKEN")

	url := "http://localhost:8086"
	client := influxdb2.NewClient(url, token)

	rand.Seed(time.Now().UnixNano())

	org := "docs"
	bucket := "go-playground"
	writeAPI := client.WriteAPIBlocking(org, bucket)
	for value := 0; value < 1000; value++ {
		tags := map[string]string{
			"tagname1": "tagvalue1",
		}
		myValue := randomInRangeInt(-100, 100)

		fields := map[string]interface{}{
			"field1":          myValue,
			"estado":          "pepe",
			"engine_speed":    randomInRangeFloat(800, 8500),
			"engine_torque":   randomInRangeFloat(100, 450),
			"oil_temperature": randomInRangeFloat(80, 110),
		}
		fmt.Printf("fields: %v\n", fields)

		point := write.NewPoint("measurement1", tags, fields, time.Now())
		time.Sleep(100 * time.Millisecond) // separate points by 1 second

		if err := writeAPI.WritePoint(context.Background(), point); err != nil {
			log.Fatal(err)
		}
	}
}
