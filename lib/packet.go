package lib

import (
//	""
)

type packet struct {
	data      []byte
	port      uint8
	channel   uint8
	payload   []byte
	writeable bool
}

func CrazyPacket(data [32]byte) *packet {
	p := &packet{}
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

/*
function Crazypacket(data)
{
if (data)
{
this.data = data;
this._port = (this.data[0] & 0xF0) >> 4;
this._channel = this.data[0] & 0x03;
this.payload = data.slice(1);
this._writable = false;
}
else
{
this.data = new Buffer(32);
this.data.fill(0);
this._writable = true;
}
this.pointer = 1;
}
*/
