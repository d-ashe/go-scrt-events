package main

import (
	"get-scrt-events-go/cmd"
)

func main() {
	if err := cmd.ScrtEventsCmd().Execute(); err != nil {
		panic(err)
	}
}
