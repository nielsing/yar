package main

import (
	"github.com/Furduhlutur/yar/robber"
)

// Entry function handles CLI arguments and executes accordingly
func main() {
	m := robber.NewMiddleware()
	m.Start()
	robber.CleanUp(m)
}
