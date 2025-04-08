package linklayerprotocol

import "math/bits"

//REFIN and REFOUT hardcoded
type Crc16 struct {
	Polynomial uint16 // Polynomial used in this CRC calculation
	StartMask  uint16 // bit shift register init
	EndMask    uint16 // xor-out, applied before returning result
	value      uint16 // value of processing steps
}

func (crc *Crc16) Init() {
	crc.value = crc.StartMask
}

func (crc *Crc16) Update(data uint8) {
	crc.value ^= uint16(data) << 8
	for i := 0; i < 8; i++ {
		leftBitIsOne := (crc.value & 0x8000) != 0
		crc.value <<= 1
		if leftBitIsOne {
			crc.value ^= crc.Polynomial
		}
	}
}

func (crc *Crc16) UpdateReflect(data uint8) {
	crc.value ^= uint16(data)
	for i := 0; i < 8; i++ {
		rightBitIsOne := (crc.value & 0x0001) != 0
		crc.value >>= 1
		if rightBitIsOne {
			crc.value ^= crc.Polynomial
		}
	}
}

func (crc *Crc16) Result() uint16 {
	//doesnt store the result in crc.value. So the end user can call several times Result()
	return crc.value ^ crc.EndMask
}

// Compute from init to result
// This method don't use Crc.Value member.
func (crc *Crc16) Compute(data []uint8) uint16 {
	value := crc.StartMask
	for i := 0; i < len(data); i++ {
		value ^= uint16(data[i]) << 8
		for j := 0; j < 8; j++ {
			leftBitIsOne := (value & 0x8000) != 0
			value <<= 1
			if leftBitIsOne {
				value ^= crc.Polynomial
			}
		}
	}
	return value ^ crc.EndMask
}

// Compute from init to result in reflect mode (REFIN and REFOUT)
// This method don't use Crc.Value member.
func (crc *Crc16) ComputeReflect(data []uint8) uint16 {
	reversePoly := bits.Reverse16(crc.Polynomial)
	value := crc.StartMask
	for i := 0; i < len(data); i++ {
		value ^= uint16(data[i])
		for j := 0; j < 8; j++ {
			rightBitIsOne := (value & 0x0001) != 0
			value >>= 1
			if rightBitIsOne {
				value ^= reversePoly
			}
		}
	}
	return value ^ crc.EndMask
}
