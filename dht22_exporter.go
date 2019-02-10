package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/morus12/dht22"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var (
		addr        = flag.String("listen-address", ":9543", "The address to listen on for HTTP requests.")
		metricsPath = flag.String("metrics-path", "/metrics", "The address to listen on for HTTP requests.")
		gpioPort    = flag.String("gpio-port", "4", "The GPIO port where DHT22 is connected.")
	)

	flag.Parse()

	if err := prometheus.Register(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Subsystem: "dht22",
			Name:      "temperature_celsius",
			Help:      "Temperature in Celsius",
		},
		func() float64 {
			sensor := dht22.New(*gpioPort)
			temperature, err := sensor.Temperature()
			if err != nil {
				log.Fatal(err)
			}
			return float64(temperature)
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
			sensor := dht22.New(*gpioPort)
			humidity, err := sensor.Humidity()
			if err != nil {
				log.Fatal(err)
			}
			return float64(humidity)
		},
	)); err == nil {
		fmt.Println("GaugeFunc 'humidity_percent', registered.")
	}

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>DHT22 Exporter</title></head>
             <body>
             <h1>DHT22 Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})
	log.Fatal(http.ListenAndServe(*addr, nil))
}
