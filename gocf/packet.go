package gocf

import (
	"bytes"
	"fmt"
	"log"
)

func Init() (*CrazyRadio, error) {
	c := &CrazyRadio{
		pacman: make(chan packet),
		pacs:   []packet{},
	}

	// Handle packets
	go c.packetHandler()

	return c, nill
}

type packet struct {
	data      []byte
	port      uint8
	channel   uint8
	payload   []byte
	writeable bool
	typ       string
}

// Create a new packet for the crazyflie
func newPacket(data []byte) *packet {
	p := &packet{}
	var b byt
	if data != [32]byte{} {
		p.data = data
		p.port = (data[0] & 0xF0) >> 4
		p.channel = data[0] & 0x03
		p.payload = data[1:]
		p.writeable = false
	} else {
		p.data = [32]byte{}
		p.writeable = true
	}
	return p
}

func (c *CrazyRadio) packetHandler() {
	ship := make(chan bool, 1)
	for {
		select {
		case pac := <-c.pacman:
			c.pacs = append(c.pacs, pac)
			ship <- true
		case <-ship:
			if len(c.pacs) > 0 {
				go c.packetManager(ship)
			} else {
				c.pacman <- c.pacman
			}
		}
	}
}

func (c *CrazyRadio) packetManager(ship chan bool) {
	pac := c.packetPop(0)
	typ := pac.typ
	for i, p := range c.pacs {
		if p.typ == typ {
			pac = c.packetPop(i)
		}
	}

	c.packetShip(pac)
	ship <- true
}

func (c *CrazyRadio) packetShip(pac packet) {
	ack, err := c.sendPacket(pac.payload)
	if err != nil {
		log.Println("Packet err: ", err)
		c.pacs = append([]packet{pac}, c.pacs...)
	}
	// Retry ...
}

func (c *CrazyRadio) packetPop(i int) packet {
	pac := c.pacs[i]
	c.pacs = append(c.pacs[:i], c.pacs[i+1]...)
	return pac
}
