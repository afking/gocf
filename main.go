package main

import (
	"github.com/afking/crazyflie/gocf"
	"log"
	"time"
)

func main() {
	// Initialise
	c, err := lib.Init()
	//defer c.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Initialised")

	// Main function commands
	for i := 0; i < 20; i++ {
		c.SetPoint(0, 0, 0, 20000)
		time.Sleep(time.Second * 0.1)
	}

	// Close ports

	// Run commands
	/*
		// Setup driver
		ctx, dev, err := driver.Drive(true) // Currently debugging
		if err != nil {
			log.Fatal("CrazyRadio driver error: ", err)
		}
		defer driver.Close(ctx, dev)

		//dev.Open
	*/
}
