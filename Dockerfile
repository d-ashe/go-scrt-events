FROM golang:latest as build

WORKDIR /go/src/get-scrt-events-go

COPY . .

RUN go get -d -v
RUN go build 

RUN chmod +x get-scrt-events-go
ENTRYPOINT [ "./go-scrt-events" ]

CMD ["--config", "config.yml", "-v", "debug"]