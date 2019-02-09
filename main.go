package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	dht "github.com/d2r2/go-dht"
	logger "github.com/d2r2/go-logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const sensorType = dht.DHT22
const gpioPort = 4

func getStats() (temp, hum float64) {
	var (
		temperature, humidity float32
		err                   error
	)

	temperature, humidity, err = dht.ReadDHTxx(sensorType, gpioPort, true)

	if err != nil {
		log.Print(err)
	}

	if humidity == 0 {
		log.Print(err)
	}

	return float64(temperature), float64(humidity)
}

func main() {
	addr := flag.String("listen-address", ":8090", "The address to listen on for HTTP requests.")

	flag.Parse()

	logger.ChangePackageLogLevel("dht", logger.ErrorLevel)

	if err := prometheus.Register(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Subsystem: "dht22",
			Name:      "temperature_celsius",
			Help:      "Temperature in Celsius",
		},
		func() float64 {
			temperature, _ := getStats()
			return temperature
		},
	)); err == nil {
		fmt.Println("GaugeFunc 'temperature_celsius', registered.")
	}

	if err := prometheus.Register(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Subsystem: "dht22",
			Name:      "humidity_percent",
			Help:      "Humidity in percent",
		},
		func() float64 {
			_, humidity := getStats()
			return humidity
		},
	)); err == nil {
		fmt.Println("GaugeFunc 'humidity_percent', registered.")
	}

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
