package lib

//import ()

// High level Commanding funcitons
func (c *CrazyRadio) SetPoint(roll, pitch, yaw float32, thrust uint16) {
	roll = 0.707 * (roll - pitch)
	pitch = 0.707 * (roll + pitch)

	p := &packet{
		port: 1,
		data: 1,
	}

	c.sendPacket(dataOut)
	/*
	   	if self._x_mode:
	               roll = 0.707 * (roll - pitch)
	               pitch = 0.707 * (roll + pitch)

	           pk = CRTPPacket()
	           pk.port = CRTPPort.COMMANDER
	           pk.data = struct.pack('<fffH', roll, -pitch, yaw, thrust)
	           self._cf.send_packet(pk)
	*/
}

/*
CrazyDriver.prototype.setpoint = function(roll, pitch, yaw, thrust)
{
var self = this,
deferred = P.defer();

var packet = new Crazypacket();
packet.port = Protocol.Ports.COMMANDER;

packet.writeFloat(roll)
.writeFloat(-pitch)
.writeFloat(yaw)
.writeUnsignedShort(thrust)
.endPacket();

return this.radio.sendPacket(packet)
.then(function(item)
{
return item;
})
.fail(function(err)
{
console.log('failure in setpoint');
console.log(err);
});
};

/ header management

Crazypacket.prototype.updateHeader = function()
{
this.data[0] = ((this._port & 0x0f) << 4 | 0x3 << 2 | (this._channel & 0x03));
};

Crazypacket.prototype.getHeader = function()
{
return this.data[0];
};

Crazypacket.prototype.setHeader = function(value)
{
this.data[0] = value;
return this;
};

Crazypacket.prototype.__defineGetter__('header', Crazypacket.prototype.getHeader);
Crazypacket.prototype.__defineSetter__('header', Crazypacket.prototype.setHeader);

Crazypacket.prototype.setChannel = function(value)
{
this._channel = value;
this.updateHeader();
return this;
};

Crazypacket.prototype.getChannel = function(value)
{
return this._channel;
};

Crazypacket.prototype.__defineGetter__('channel', Crazypacket.prototype.getChannel);
Crazypacket.prototype.__defineSetter__('channel', Crazypacket.prototype.setChannel);


Crazypacket.prototype.setPort = function(value)
{
this._port = value;
this.updateHeader();
return this;
};

Crazypacket.prototype.getPort = function(value)
{
return this._port;
};

Crazypacket.prototype.__defineGetter__('port', Crazypacket.prototype.getPort);
Crazypacket.prototype.__defineSetter__('port', Crazypacket.prototype.setPort);


// --------------------------------------------

function RadioAck(buffer)
{
this.ack = false;
this.powerDet = false;
this.retry = 0;
this.data = undefined;

if (buffer)
this.parse(buffer);
}

RadioAck.prototype.parse = function(buffer)
{
this.ack = (buffer[0] & 0x01) !== 0;
this.powerDet = (buffer[0] & 0x02) !== 0;
this.retry = buffer[0] >> 4;
this.data = buffer.slice(1);
this.packet = new Crazypacket(this.data);
};

// --------------------------------------------
*/
