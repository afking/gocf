// lsusb lists attached USB devices.
package gocf

import (
	_ "encoding/binary"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/kylelemons/gousb/usb"
	"github.com/kylelemons/gousb/usbid"
)

var (
	// USB Paramteres
	CRADIO_VID = 0x1915
	CRADIO_PID = 0x7777
	// a libusb constant
	TYPE_VENDOR               uint8 = 0x02 << 5
	LIBUSB_REQUEST_GET_STATUS uint8 = 0x0
	LIBUSB_ENDPOINT_IN        uint8 = 0x80
	LIBUSB_ENDPOINT_OUT       uint8 = 0x00
	// Radio commands
	SET_RADIO_CHANNEL uint8 = 0x01
	SET_RADIO_ADDRESS uint8 = 0x02
	SET_DATA_RATE     uint8 = 0x03
	SET_RADIO_POWER   uint8 = 0x04
	SET_RADIO_ARD     uint8 = 0x05
	SET_RADIO_ARC     uint8 = 0x06
	ACK_ENABLE        uint8 = 0x10
	SET_CONT_CARRIER  uint8 = 0x20
	SCAN_CHANNELS     uint8 = 0x21
	LAUNCH_BOOTLOADER uint8 = 0xFF
	// Data Rates
	DR_250KPS uint16 = 0
	DR_1MPS   uint16 = 1
	DR_2MPS   uint16 = 2
	// Power Rates
	P_M18DBM uint16 = 0
	P_M12DBM uint16 = 1
	P_M6DBM  uint16 = 2
	P_0DBM   uint16 = 3

	debug = flag.Int("debug", 0, "libusb debug level (0..3)")
)

type CrazyRadio struct {
	ctx      *usb.Context
	dev      *usb.Device
	in       usb.Endpoint
	out      usb.Endpoint
	version  uint16
	arc      int16
	address  []byte
	channel  uint16
	datarate uint16
	pacman   chan *packet
	pacs     []*packet
}

func (c *CrazyRadio) CrazyRadioDrive() error {
	// One context should be opened for the application.
	c.ctx = usb.NewContext()
	c.ctx.Debug(*debug)

	// ListDevices is used to find the devices to open.
	devs, err := c.ctx.ListDevices(func(desc *usb.Descriptor) bool {
		if desc.Vendor == usb.ID(CRADIO_VID) && desc.Product == usb.ID(CRADIO_PID) {
			// The usbid package can be used to print out human readable information.
			fmt.Printf("%03d.%03d %s:%s %s\n", desc.Bus, desc.Address, desc.Vendor, desc.Product, usbid.Describe(desc))
			fmt.Printf(" Protocol: %s\n", usbid.Classify(desc))

			for _, cfg := range desc.Configs {
				fmt.Printf(" %s:\n", cfg)
				for _, alt := range cfg.Interfaces {
					fmt.Printf(" --------------\n")
					for _, iface := range alt.Setups {
						fmt.Printf(" %s\n", iface)
						fmt.Printf(" %s\n", usbid.Classify(iface))
						for _, end := range iface.Endpoints {
							fmt.Printf(" %s\n", end)
						}
					}
				}
				fmt.Printf(" --------------\n")
			}
			fmt.Println("Returning True")
			return true
		} else {
			return false
		}
	})
	if err != nil {
		return err
	} else if len(devs) == 0 {
		err = fmt.Errorf("CrazyRadio was not detected.")
		return err
	}

	// Seetup Device, Grab first radio
	c.dev = devs[0]

	// Sets timeout configurations
	setTimeouts(c.dev)
	c.datarate = 2
	c.channel = 2
	c.arc = -1

	c.in, err = c.dev.OpenEndpoint(01, 00, 00, 01|uint8(usb.ENDPOINT_DIR_IN))
	if err != nil {
		return err
	}

	c.out, err = c.dev.OpenEndpoint(01, 00, 00, 01|uint8(usb.ENDPOINT_DIR_OUT))
	if err != nil {
		return err
	}

	log.Println(c.out.Info().PollInterval)

	c.reset()

	return nil
}

func (c *CrazyRadio) reset() {
	var i uint8 = 0xE7
	c.address = []byte{i, i, i, i, i}
	//var i uint8 = 0xE7
	//c.address = make([]byte, 5)
	//binary.LittleEndian.PutUint16(c.address, i)

	c.arc = -1
	c.setChannel(c.channel)
	c.setDataRate(c.datarate)
	c.setContCarrier(false)
	c.setAddress(c.address)
	c.setPower(P_0DBM)
	c.setAckRetryCount(3)
	c.setArdBytes(32)
}

func setTimeouts(d *usb.Device) {
	d.ReadTimeout = 1000 * time.Millisecond
	d.WriteTimeout = 1000 * time.Millisecond
	d.ControlTimeout = 1000 * time.Millisecond
}

func (c *CrazyRadio) sendVendorSetup(request uint8, value, index uint16, data []byte) error {
	n, err := c.dev.Control(TYPE_VENDOR, request, value, index, data)
	log.Println("sendVendorSetup: ", n)
	if err != nil {
		log.Println("sendVendorSetup error:", err)
		return err
	}
	return nil
}

func (c *CrazyRadio) receiveVendor(request uint8, value, index uint16, length int) ([]byte, error) {
	b := make([]byte, length)
	n, err := c.dev.Control(TYPE_VENDOR|LIBUSB_ENDPOINT_IN, request, value, index, b)
	log.Println("recieveVendorSetup: ", n)
	if err != nil {
		log.Println("receiveVendor error:", err)
	}
	return b, err
}

// Dongle Configurations
func (c *CrazyRadio) setChannel(channel uint16) {
	c.sendVendorSetup(SET_RADIO_CHANNEL, channel, 0, []byte{})
}
func (c *CrazyRadio) setAddress(address []byte) {
	c.sendVendorSetup(SET_RADIO_ADDRESS, 0, 0, address)
}
func (c *CrazyRadio) setDataRate(rate uint16) {
	c.sendVendorSetup(SET_DATA_RATE, rate, 0, []byte{})
}
func (c *CrazyRadio) setPower(level uint16) {
	c.sendVendorSetup(SET_RADIO_POWER, level, 0, []byte{})
}
func (c *CrazyRadio) setAckRetryCount(retries uint16) {
	c.sendVendorSetup(SET_RADIO_ARC, retries, 0, []byte{})
}
func (c *CrazyRadio) setArdRetryDelay(delay uint16) {
	t := int((delay / 250) - 1)
	if t < 0 {
		t = 0
	} else if t > 0xF {
		t = 0xF
	}
	c.sendVendorSetup(SET_RADIO_ARD, uint16(t), 0, []byte{})
}
func (c *CrazyRadio) setArdBytes(nbytes uint16) {
	c.sendVendorSetup(SET_RADIO_ARD, 0x80|nbytes, 0, []byte{})
}
func (c *CrazyRadio) setContCarrier(active bool) {
	x := 0
	if active {
		x = 1
	}
	c.sendVendorSetup(SET_CONT_CARRIER, uint16(x), 0, []byte{})
}
func (c *CrazyRadio) scanChannels(start, stop uint16, packet []byte) ([]byte, error) {
	//var channels []int
	err := c.sendVendorSetup(SCAN_CHANNELS, start, stop, packet)
	if err != nil {
		return []byte{}, err
	}
	b, err := c.receiveVendor(SCAN_CHANNELS, 0, 0, 64)
	fmt.Println("scanChannels")
	return b, err
}

// Data Transfer
func (c *CrazyRadio) sendPacket(dataOut []byte) (Ack, error) {
	var ackIn Ack
	var data []byte

	_, err := c.out.Write(dataOut)
	if err != nil {
		ackIn.ack = false
		return ackIn, err
	}

	n, err := c.in.Read(data)
	if err != nil {
		ackIn.ack = false
		return ackIn, err
	}

	if n == 0 {
		return ackIn, nil
	}

	if ackIn.ack = false; (data[0] & 0x01) != 0 {
		ackIn.ack = true
	}
	if ackIn.powerDet = false; (data[0] & 0x02) != 0 {
		ackIn.powerDet = true
	}
	ackIn.retry = data[0] >> 4
	ackIn.data = data[1:]

	return ackIn, nil
}

func (c *CrazyRadio) Write(data []byte) error {
	_, err := c.out.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// Close currently opened ports
func (c *CrazyRadio) Close() {
	c.dev.Close()
	c.ctx.Close()
}
