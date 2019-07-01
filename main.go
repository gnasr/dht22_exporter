package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/morus12/dht22"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	listenAddress = flag.String("listen-address", ":9543", "The address to listen on for HTTP requests.")
	metricsPath   = flag.String("metrics-path", "/metrics", "The path of the metrics endpoint.")
	gpioPort      = flag.String("gpio-port", "4", "The GPIO port where DHT22 is connected.")
)

func main() {
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
				log.Fatalf("error reading temperature %v", err)
			}

			return float64(temperature)
		},
	)); err != nil {
		log.Fatalf("error registering gaugefunc: %v", err)
	}
	log.Println("GaugeFunc 'temperature_celsius', registered.")

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
				log.Fatalf("error reading humidity %v", err)
			}

			return float64(humidity)
		},
	)); err != nil {
		log.Fatalf("error registering gaugefunc: %v", err)
	}
	log.Println("GaugeFunc 'temperature_celsius', registered.")

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
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
