FROM golang:latest as build

WORKDIR /go/src/go-scrt-events

COPY . .

RUN go get -d -v
RUN go build 

RUN chmod +x go-scrt-events
ENTRYPOINT [ "./go-scrt-events" ]

CMD ["--config", "config.yml", "-v", "debug"]