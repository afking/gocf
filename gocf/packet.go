package gocf

import (
	"log"
)

type packet struct {
	data      []byte
	port      uint8
	channel   uint8
	payload   []byte
	writeable bool
	typ       string
}

type Ack struct {
	ack      bool
	powerDet bool
	retry    byte
	data     []byte
}

// Create a new packet for the crazyflie
func newPacket(data []byte, typ string) *packet {
	p := &packet{}
	//var b byte
	if data != nil {
		p.data = data
		p.port = (data[0] & 0xF0) >> 4
		p.channel = data[0] & 0x03
		p.payload = data[1:]
		p.writeable = false
		p.typ = typ + string(p.channel)
	} else {
		p.data = []byte{}
		p.writeable = true
	}
	return p
}

func (c *CrazyRadio) packetHandler() {
	ship := make(chan bool, 1)
	pacs := []*packet{}

	// Init
	ship <- true

	for {
		select {
		case pac := <-c.pacman:
			pacs = append(pacs, pac)
		case <-ship:
			// Testing for buffered packets
			if len(pacs) > 1 {
				log.Println("Mulitple Packets buffered")
			}

			// Append handler packets
			c.pacs = append(c.pacs, pacs...)
			pacs = []*packet{}

			if len(c.pacs) > 0 {
				go c.packetManager(ship)
			} else {
				pac := <-c.pacman // Waits
				pacs = append(pacs, pac)
				ship <- true
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

func (c *CrazyRadio) packetShip(pac *packet) {
	//ack, err := c.sendPacket(pac.payload)
	_, err := c.sendPacket(pac.payload)
	if err != nil {
		log.Println("Packet err: ", err)
		//c.pacs = append([]packet{pac}, c.pacs...)
	}
	// Retry ...
}

func (c *CrazyRadio) packetPop(i int) *packet {
	pac := c.pacs[i]
	c.pacs = append(c.pacs[:i], c.pacs[i+1:]...)
	return pac
}
