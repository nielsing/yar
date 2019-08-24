package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Furduhlutur/yar/robber"
)

func main() {
	m := robber.NewMiddleware()

	kill := make(chan bool)
	cleanup := make(chan bool)
	finished := make(chan bool)
	sigc := make(chan os.Signal)
	signal.Notify(sigc, syscall.SIGINT)
	go robber.HandleSigInt(m, sigc, kill, finished, cleanup)

	m.Start(kill, finished, cleanup)
}
