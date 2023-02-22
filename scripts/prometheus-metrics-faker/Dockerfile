FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /prometheus-metrics-faker

EXPOSE 9000

ENTRYPOINT [ "/prometheus-metrics-faker" ]
