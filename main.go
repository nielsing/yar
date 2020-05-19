package main

import (
	"os"
	"os/signal"

	"github.com/nielsing/yar/robber"
)

func main() {
	m := robber.NewMiddleware()

	kill := make(chan bool)
	cleanup := make(chan bool)
	finished := make(chan bool)
	sigc := make(chan os.Signal, 2)
	signal.Notify(sigc, os.Interrupt)
	go robber.HandleSigInt(m, sigc, kill, finished, cleanup)
	go robber.HandleSigInt(m, sigc, kill, finished, cleanup)

	m.Start(kill, finished, cleanup)
	if m.Flags.SavePresent {
		robber.SaveFindings(m)
	}
}
