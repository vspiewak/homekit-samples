package main

import (
	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"

	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// Create the switch accessory
	a := accessory.NewSwitch(accessory.Info{
		Name: "My Golang Switch",
	})

	// Log switch state on update
	a.Switch.On.OnValueRemoteUpdate(func(on bool) {
		if on {
			log.Println("Switch is ON")
		} else {
			log.Println("Switch is OFF")
		}
	})

	// Store the data in the "./db" directory
	fs := hap.NewFsStore("./db")

	// Create the hap server
	server, err := hap.NewServer(fs, a.A)
	if err != nil {
		// stop if an error happens
		log.Panic(err)
	}

	// Set HomeKit Pin
	server.Pin = "00112233"

	// Listen for interrupts and SIGTERM signals to stop the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-c
		// Stop delivering signals
		signal.Stop(c)
		// Cancel the context to stop the server
		cancel()
	}()

	// Run the server
	server.ListenAndServe(ctx)

}
