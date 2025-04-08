package linklayerprotocol

const (
	StartByte      uint8 = 0x7E
	MaxPayloadSize uint8 = 10
	OverheadSize   uint8 = 4 // start, length, crc16
	MaxFrameSize   uint8 = MaxPayloadSize + OverheadSize
)

type LinkProtocol struct {
	isUnframing bool
	rxLength    uint8
	rxBuffer    [MaxFrameSize]uint8
	rxPosition  uint8
	lowerLayer  ComSender
	upperLayer  AppReceiver
	crc         Crc16
	crcBytes    uint16
}

func NewDefaultLinkProtocol(lowerLayer ComSender, upperLayer AppReceiver, crc Crc16) LinkProtocol {
	return LinkProtocol{
		lowerLayer: lowerLayer,
		upperLayer: upperLayer,
		crc:        crc,
	}
}

func (lp *LinkProtocol) Init(lowerLayer ComSender, upperLayer AppReceiver, crc Crc16) {
	lp.lowerLayer = lowerLayer
	lp.upperLayer = upperLayer
	lp.crc = crc
}

func (lp *LinkProtocol) reset() {
	lp.isUnframing = false
	//lp.rxLength = 0
	lp.rxPosition = 0
	lp.crcBytes = 0
}

func (lp *LinkProtocol) ReceiveByte(currentByte uint8) {
	// wait start byte
	if !lp.isUnframing {
		if currentByte == StartByte {
			lp.isUnframing = true
			lp.rxLength = 0
			return
		} else {
			return
		}
	}
	// read length byte
	if lp.rxLength == 0 {
		lp.rxLength = currentByte
		lp.crc.Init()
		return
	}
	// just fill the buffer and update crc
	if lp.rxPosition < lp.rxLength {
		lp.rxBuffer[lp.rxPosition] = currentByte
		lp.rxPosition++
		lp.crc.Update(currentByte)
		return
	}
	// save first byte of CRC16
	if lp.rxPosition <= lp.rxLength {
		lp.rxPosition++
		lp.crcBytes = uint16(currentByte) << 8
		return
	}
	// build CRC16 and compare
	lp.crcBytes |= uint16(currentByte)
	if lp.crcBytes != lp.crc.Result() {
		// silent drop this faulty frame.
		// TODO event or counter ???
		lp.reset()
		return
	}
	// all good, transmit to upper layer
	lp.upperLayer.Receive(lp.rxBuffer[:lp.rxLength])
	lp.reset()
}

func (lp *LinkProtocol) Send(bytes []uint8) {
	size := uint8(len(bytes))
	if size > MaxPayloadSize {
		return
	}
	frameSize := size + OverheadSize
	var frame []uint8 = make([]uint8, frameSize)
	frame[0] = StartByte
	frame[1] = size

	var crc uint16
	for index, cByte := range bytes {
		frame[index+2] = cByte
		crc = lp.crc.Compute(bytes)
	}

	frame[frameSize-2] = uint8(crc >> 8)
	frame[frameSize-1] = uint8(crc)

	lp.lowerLayer.Send(frame)
}
