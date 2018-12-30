package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	dht "github.com/d2r2/go-dht"
	logger "github.com/d2r2/go-logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr                  string
	gpioPort              int
	stype                 string
	boostPerfFlag         bool
	prometheusTemperature = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "dht22_temperature_celsius",
			Help: "DHT22 Temperature in celsius currently collected",
		})

	prometheusHumidity = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "dht22_humidity_percentage",
			Help: "DHT22 Humidity percentage currently collected",
		})
)

func fetchDataFromDHT() {
	temperature, humidity, retried, err :=
		dht.ReadDHTxxWithRetry(getSensorType(), gpioPort, boostPerfFlag, 10)
	if err != nil {
		log.Fatal(err)
	}

	prometheusTemperature.Set(float64(temperature))
	prometheusHumidity.Set(float64(humidity))
	// Print temperature and humidity
	fmt.Printf("Temperature = %v*C, Humidity = %v%% (retried %d times)\n",
		temperature, humidity, retried)
}

func init() {
	flag.StringVar(&addr, "listen-address", ":8080", "The address to listen on for HTTP requests.")
	flag.IntVar(&gpioPort, "gpio-port", 4, "GPIO Port")
	flag.StringVar(&stype, "sensor-type", "dht22", "sensor type (dht22, dht11)")
	flag.BoolVar(&boostPerfFlag, "boost", false, "boost performance")
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(prometheusTemperature)
	prometheus.MustRegister(prometheusHumidity)
}

func getSensorType() dht.SensorType {
	var sensorType dht.SensorType

	if stype == "dht22" || stype == "am2302" {
		sensorType = dht.DHT22
	} else if stype == "dht11" {
		sensorType = dht.DHT11
	}

	return sensorType
}

func main() {
	flag.Parse()
	logger.ChangePackageLogLevel("dht", logger.ErrorLevel)

	ticker := time.NewTicker(2 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				fetchDataFromDHT()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(addr, nil))
}
