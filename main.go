package main

import (
	"github.com/secretanalytics/go-scrt-events/cmd"
)

func main() {
	if err := cmd.ScrtEventsCmd().Execute(); err != nil {
		panic(err)
	}
}
