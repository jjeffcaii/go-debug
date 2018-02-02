# go-debug

a colorful debug tool in golang.

## Install

```bash
$ dep ensure -add github.com/jjeffcaii/go-debug
```

## Example

```go
package main

import (
	"github.com/jjeffcaii/go-debug"
)

func main() {
	// env settings:
	// show all: DEBUG=*
	// show service only: DEBUG=service:*
	// show http and redis: DEBUG=http,redis
	debug.Debug("http").Printf("this is a test debug info: %s\n", "hello world")
	debug.Debug("service:user", debug.UpperCase).Println("this is a test debug info")
	debug.Debug("service:bill", debug.UpperCase).Println("this is a test debug info")
	debug.Debug("mysql", debug.TimeLocal).Println("this is a test debug info")
	debug.Debug("redis", debug.TimeUTC|debug.UpperCase).Println("this is a test debug info")
}

```

![screen_shot](screen_shot.png "screen_shot.png")
