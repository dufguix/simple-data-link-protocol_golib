# Data Link Layer Protocol

Protocol for simple simplex/duplex communication.
- CRC to check data integrity
- drop frame in silent mode when crc is bad
- build/unbuild frame byte after byte
- transfer payload to another party

Frame: Start / Payload_Length / Payload / CRC16

Communication layer between bits and upper layer.

This package build and unbuild frames.
For communication transfer, use another package for serial_port on computer or uart on microcontroller (TinyGo).

Thirdparty can easily call LinkProtocol.SendByte().
LinkProtocol call ComSender.send() to send bytes over serial communication.

Thirparty can easily call LinkProtocol.ReceiveByte() when serial com has new data.
LinkProtocol call AppReceiver.receive() interface to send payload to the upper layer.

C lib available [here](https://github.com/dufguix/simple-data-link-protocol_clib).

## Install
```
go get github.com/dufguix/simple-data-link-protocol_golib

import (
	"time"
	sdlp "github.com/dufguix/simple-data-link-protocol_golib"
)
```

## TODO
- provide examples in readme