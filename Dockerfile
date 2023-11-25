FROM golang:1.21

WORKDIR /app

COPY go.mod *.go ./

RUN CGO_ENABLED=0 go build -o /exchange-rates-api

CMD ["/exchange-rates-api"]
