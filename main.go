package main

import (
	"github.com/kctjohnson/bubble-boids/cmd/cli"
	"github.com/kctjohnson/bubble-boids/cmd/server"
)

const servermode = false

func main() {
	if servermode {
		server.Execute()
	} else {
		cli.Execute()
	}
}
