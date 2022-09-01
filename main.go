package main

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/kctjohnson/bubble-boids/cmd/cli"
	"github.com/kctjohnson/bubble-boids/cmd/server"
)

const servermode = false

func main() {
	go func() {
		http.ListenAndServe(":8080", nil)
	}()
	if servermode {
		server.Execute()
	} else {
		cli.Execute()
	}
}
