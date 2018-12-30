# DHT22 exporter

This is a simple DHT22 exporter that uses internally https://github.com/d2r2/go-dht. It will expose the temperature and humidity as Prometheus metrics:

```bash
$ curl -s 192.168.178.48:8080/metrics | grep dht
# HELP dht22_humidity_percentage DHT22 Humidity percentage currently collected
# TYPE dht22_humidity_percentage gauge
dht22_humidity_percentage 49.19999694824219
# HELP dht22_temperature_celsius DHT22 Temperature in celsius currently collected
# TYPE dht22_temperature_celsius gauge
dht22_temperature_celsius 18.0
```
