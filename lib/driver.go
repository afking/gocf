// lsusb lists attached USB devices.
package lib

import (
	"fmt"
	_ "log"
	"time"

	"github.com/kylelemons/gousb/usb"
	"github.com/kylelemons/gousb/usbid"
)

var (
	// USB Paramteres
	CRADIO_VID = 0x1915
	CRADIO_PID = 0x7777
	// a libusb constant
	TYPE_VENDOR               = 0x02 << 5
	LIBUSB_REQUEST_GET_STATUS = 0x0
	// Radio commands
	SET_RADIO_CHANNEL = 0x01
	SET_RADIO_ADDRESS = 0x02
	SET_DATA_RATE     = 0x03
	SET_RADIO_POWER   = 0x04
	SET_RADIO_ARD     = 0x05
	SET_RADIO_ARC     = 0x06
	ACK_ENABLE        = 0x10
	SET_CONT_CARRIER  = 0x20
	SCANN_CHANNELS    = 0x21
	LAUNCH_BOOTLOADER = 0xFF
	// Data Rates
	DR_250KPS = 0
	DR_1MPS   = 1
	DR_2MPS   = 2
	// Power Rates
	P_M18DBM = 0
	P_M12DBM = 1
	P_M6DBM  = 2
	P_0DBM   = 3
)

type CrazyRadio struct {
	device      *usb.Device
	inEndpoint  *usb.Endpoint
	outEndpoint *usb.Endpoint
	version     int16
	arc         int16
	address     int16
	//outstream int16
	//instream  int16
	channel  int16
	datarate int16
	//pingTimer    int16
	//pingInterval int16
}

func CrazyRadioDrive() error {
	devices, err := findCrazyRadios()
	if len(devices) == 0 {
		err = fmt.Errorf("CrazyRadio was not detected")
	}
	if err != nil {
		return err
	}
	dev := devices[0]
	defer dev.Close()

	setTimeouts(dev) // Sets timeout configurations

	c := &CrazyRadio{
		channel:  2,
		datarate: 2,
		device:   dev,
	}

	// This is most likely wrong
	c.inEndpoint, err = dev.OpenEndpoint(0, 0, 0, 0) // conf, iface, setup, epoint uint8
	if err != nil {
		return err
	}
	c.outEndpoint, err = dev.OpenEndpoint(0, 0, 0, 1) // conf, iface, setup, epoint uint8
	if err != nil {
		return err
	}

	c.reset()

	return nil
}

func (c *CrazyRadio) reset() {
	c.arc = -1
	c.address = [5]byte{0xE7}
	c.setChannel(c.channel)
	c.setDataRate(c.datarate)
	c.setContCarrier(false)
	c.setAddress(c.address)
	c.setPower(P_0DBM)
	c.setAckRetryCount(3)
	c.setArdBytes(32)
}

func findCrazyRadios() ([]*usb.Device, error) {
	// Only one context should be needed for an application. It should always be closed.
	ctx := usb.NewContext()
	defer ctx.Close()

	//ctx.Debug(*debug)

	// ListDevices is used to find the devices to open.
	devs, err := ctx.ListDevices(func(desc *usb.Descriptor) bool {
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

		if desc.Vendor == usb.ID(CRADIO_VID) && desc.Product == usb.ID(CRADIO_PID) {
			return true
		}
		return false
	})
	if err != nil {
		return nil, err
	}

	return devs, nil
}

func setTimeouts(d *usb.Device) {
	d.Readtimeout = 1000 * time.Millisecond
	d.WriteTimeout = 1000 * time.Millisecond
	d.ControlTimeout = 1000 * time.Millisecond
}

func (c *CrazyRadio) sendVendorSetup(request uint8, value, index uint16, data []byte) error {
	n, err := c.device.Control(TYPE_VENDOR, request, value, index, data)
	if err != nil {
		return err
	}
	return nil
}

// Dongle Configurations
func (c *CrazyRadio) setChannel(channel uint16) {
	c.sendVendorSetup(SET_RADIO_CHANNEL, channel, 0, []byte{})
}
func (c *CrazyRadio) setAddress(address uint16) {
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
	c.sendVendorSetup(SET_RADIO_ARD, t, 0, []byte{})
}
func (c *CrazyRadio) setArdBytes(nbytes uint16) {
	c.sendVendorSetup(SET_RADIO_ARD, 0x80|nbytes, 0, []byte{})
}
func (c *CrazyRadio) setContCarrier(active bool) {
	x := 0
	if active {
		x = 1
	}
	c.sendVendorSetup(SET_CONT_CARRIER, x, 0, []byte{})
}
func (c *CrazyRadio) scanChannels(start, stop int, packet []byte) []int {
	var channels []int
	for i := range stop - start + 1 {
		c.setChannel(i)
		status := c.sendPacket(packet)
		if status && status.ack {
			channels = append(channels, i)
		}
	}
	return channels
}

type Ack struct {
	ack      bool
	powerDet bool
	retry    int
	data     []byte
}

// Data Transfer
func (c *CrazyRadio) sendPacket(dataOut []byte) (Ack, error) {
	var ackIn Ack
	var data [64]byte

	_, err := c.endpoint.Write(dataOut)
	if err != nil {
		return nil, err
	}

	n, err := c.endpoint.Read(data)
	if err != nil {
		return nil, err
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

/*
def _send_vendor_setup(handle, request, value, index, data):
    if pyusb1:
        handle.ctrl_transfer(usb.TYPE_VENDOR, request, wValue=value,
                             wIndex=index, timeout=1000, data_or_wLength=data)
    else:
        handle.controlMsg(usb.TYPE_VENDOR, request, data, value=value,
                          index=index, timeout=1000)


// Close currently opened ports
func Close(ports ...interface{}) {
	for _, port := range ports {
		port.Close()
	}
} */

/*
// Send packets and recieve the ack fromt the radio dongle.
// Ack containts information about the packet transmition.
func SendPackets(dev *usb.Device, dataOut []byte) *Ack {
	ackIn := &Ack{}
	var data []byte

	return ackIn
} */
