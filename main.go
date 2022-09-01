package main

import (
	// "github.com/kctjohnson/bubble-boids/cmd/cli"
	"github.com/kctjohnson/bubble-boids/cmd/pixel"
	// "github.com/kctjohnson/bubble-boids/cmd/server"
)

const servermode = false

func main() {
	pixel.Execute()
	// if servermode {
	// 	server.Execute()
	// } else {
	// 	cli.Execute()
	// }
}
