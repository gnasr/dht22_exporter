FROM golang:1.11.4 AS builder

WORKDIR /build
ADD . /build
RUN go get github.com/morus12/dht22 && \
    go get github.com/prometheus/client_golang/prometheus
RUN GOOS=linux GOARCH=arm GOARM=5 go build -o dht22_exporter .

FROM  busybox:latest
LABEL maintainer="Gabriel Nasr <nasr.gab@gmail.com>"

COPY --from=builder /build/dht22_exporter /bin/dht22_exporter
ENTRYPOINT ["/bin/dht22_exporter"]
EXPOSE 9543