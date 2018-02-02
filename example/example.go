package main

import (
	"github.com/jjeffcaii/go-debug"
)

func main() {
	// env settings:
	// show all: DEBUG=*
	// show service only: DEBUG=service:*
	// show http and redis: DEBUG=http,redis
	debug.Debug("http", debug.UpperCase).Println("this is a test debug info")
	debug.Debug("service:user", debug.UpperCase).Println("this is a test debug info")
	debug.Debug("service:bill", debug.UpperCase).Println("this is a test debug info")
	debug.Debug("mysql", debug.TimeLocal).Println("this is a test debug info")
	debug.Debug("redis", debug.TimeUTC|debug.UpperCase).Println("this is a test debug info")
}
