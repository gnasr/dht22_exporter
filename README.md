# DHT22 Prometheus exporter
[DHT22](https://www.sparkfun.com/datasheets/Sensors/Temperature/DHT22.pdf) temperature and humidity sensor Prometheus exporter
## Instalation
```
go get github.com/gnasr/dht22_exporter
cd $GOPATH/src/github.com/gnasr/dht22_exporter
GOOS=linux GOARCH=arm GOARM=5 go build -o dht22_exporter
```