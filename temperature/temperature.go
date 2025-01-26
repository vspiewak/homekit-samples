package main

import (
	"time"

	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
	"github.com/shirou/gopsutil/v4/sensors"

	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// Create the temperature sensors accessory
	a := accessory.NewTemperatureSensor(accessory.Info{
		Name: "M3 Battery",
	})

	// Update loop
	go func() {
		for {

			stats, err := sensors.SensorsTemperatures()
			if err != nil {
				log.Println("error reading sensors:", err)
			} else {
				for _, stat := range stats {
					if stat.SensorKey == "NAND CH0 temp" {
						value := stat.Temperature
						log.Println("current temp:", value)
						a.TempSensor.CurrentTemperature.SetValue(value)
					}
				}
			}

			time.Sleep(5 * time.Second)

		}
	}()

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
