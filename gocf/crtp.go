package gocf

import (
	"bytes"
	"encoding/binary"
	"log"
)

func Init() *CrazyRadio {
	c := &CrazyRadio{
		pacman: make(chan *packet),
		pacs:   []*packet{},
	}

	// Start radio
	err := c.CrazyRadioDrive()
	if err != nil {
		log.Fatal(err)
	}

	// Find Crazyflie
	cfs := c.scanAtRate()

	// Connect
	c.setChannel(cfs[0])
	log.Printf("Connected to channel: %d", cfs[0])

	// Handle packets
	go c.packetHandler()

	return c
}

type CrazyDriver struct {
	uri      string
	channel  uint16
	datarate uint16
	radio    CrazyRadio
}

// Just scan at 2MBPS
func (c *CrazyRadio) scanAtRate() []uint16 {
	cfs := []uint16{}
	c.setDataRate(2) // Default atm
	var b = []uint8{0xFF}
	n, err := c.scanChannels(0, 127, b)
	if err != nil {
		log.Println(n)
		log.Fatal("scanAtRate error: ", err)
	}
	for i := 0; i < len(n); i++ {
		if n[i] != 0 {
			log.Printf("radio://1/%d/2MBPS", n[i])
			cfs = append(cfs, uint16(n[i]))
		}
	}
	if len(cfs) == 0 {
		log.Fatal("No Crazyflies Detected")
	}
	return cfs
}

// High level Commanding funcitons
func (c *CrazyRadio) SetPoint(roll, pitch, yaw float32, thrust uint16) {
	roll = 0.707 * (roll - pitch)
	pitch = 0.707 * (roll + pitch)

	buf := new(bytes.Buffer)
	var data = []interface{}{
		roll,
		pitch,
		yaw,
		thrust,
	}
	for _, v := range data {
		err := binary.Write(buf, binary.LittleEndian, v)
		if err != nil {
			log.Println("binary.Write failed: ", err)
		}
	}

	pac := newPacket(buf.Bytes(), "SetPoint")

	// Sends packet to packetHandler on channel pacman
	c.pacman <- pac
}
