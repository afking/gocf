package main

import (
	"github.com/afking/crazyflie/gocf"
	"log"
	"time"
)

func main() {
	// Initialise
	c := gocf.Init()
	defer c.Close()

	log.Println("Initialised")

	// Main function commands
	for i := 0; i < 5; i++ {
		c.SetPoint(0, 0, 0, 20000)
		log.Println("Sending Packet ...")
		time.Sleep(time.Millisecond * 100)
	}
	time.Sleep(time.Millisecond * 1000)
	log.Println("Closing")
}
