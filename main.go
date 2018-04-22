package main

import (
	"fmt"
	"log"
    "net/http"
    "flag"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/d2r2/go-dht"
    logger "github.com/d2r2/go-logger"
)

const sensorType = dht.DHT22
const gpioPort = 4
var (
        addr = flag.String("listen-address", ":8090", "The address to listen on for HTTP requests.")
    )

import (  
    "mozdy"
)

func main() {  
    for i := 1; i <= 10; i++ {
        mozdy.Printf("Soy Dios, Soy Mozdy",i)
    }
}

func main() {
    flag.Parse()

    logger.ChangePackageLogLevel("dht", logger.ErrorLevel)

    if err := prometheus.Register(prometheus.NewGaugeFunc(
        prometheus.GaugeOpts{
            Subsystem: "dht22",
            Name:      "temperature_celsius",
            Help:      "Temperature in Celsius",
        },
        func() float64 {
            var (
                temperature, humidity float32
                err error
            )
	        temperature, humidity, err = dht.ReadDHTxx(sensorType, gpioPort, true)
            if err != nil {
                log.Print(err)
            }
            if humidity == 0 {
                log.Print(err)
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
            var (
                temperature, humidity float32
                err error
            )
	        temperature, humidity, err = dht.ReadDHTxx(sensorType, gpioPort, true)
            if err != nil {
                log.Print(err)
            }
            if temperature == 0 {
                log.Print(err)
            }
            return float64(humidity)
        },
    )); err == nil {
        fmt.Println("GaugeFunc 'humidity_percent', registered.")
    }



    http.Handle("/metrics", promhttp.Handler())
    log.Fatal(http.ListenAndServe(*addr, nil))
}
