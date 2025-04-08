package linklayerprotocol

import (
	"math/bits"
	"testing"
)

func TestBitsReverse(t *testing.T) {
	if bits.Reverse16(0xa2eb) != 0xd745 {
		t.Fatalf(`bits reverse doesnt work as expected.`)
	}
}

func TestComputeFunc(t *testing.T) {
	//123456789123456789 in hex format
	input := []uint8{0x12, 0x34, 0x56, 0x78, 0x91, 0x23, 0x45, 0x67, 0x89}

	crcXmodem := Crc16{
		Polynomial: 0x1021,
		StartMask:  0,
		EndMask:    0,
	}
	xmodemResult := uint16(0x9B4C)
	xmodemComputeResult := crcXmodem.Compute(input)

	crcDds110 := Crc16{
		Polynomial: 0x8005,
		StartMask:  0x800D,
		EndMask:    0x0000,
	}
	dds110Result := uint16(0x8ACB)
	dds110ComputeResult := crcDds110.Compute(input)

	crcGsm := Crc16{
		Polynomial: 0x1021,
		StartMask:  0,
		EndMask:    0xFFFF,
	}
	gsmResult := uint16(0x64B3)
	gsmComputeResult := crcGsm.Compute(input)

	if xmodemComputeResult != xmodemResult {
		t.Fatalf(`xmodem crc: %X doesnt match %X`, xmodemComputeResult, xmodemResult)
	}
	if dds110ComputeResult != dds110Result {
		t.Fatalf(`dds110 crc: %X doesnt match %X`, dds110ComputeResult, dds110Result)
	}
	if gsmComputeResult != gsmResult {
		t.Fatalf(`gsm crc: %X doesnt match %X`, gsmComputeResult, gsmResult)
	}
}

func TestComputeReflectFunc(t *testing.T) {
	//123456789123456789 in hex format
	input := []uint8{0x12, 0x34, 0x56, 0x78, 0x91, 0x23, 0x45, 0x67, 0x89}

	crcCcitt := Crc16{
		Polynomial: 0x1021,
		StartMask:  0,
		EndMask:    0,
	}

	ccittResult := uint16(0xC0F3)
	ccittComputeResult := crcCcitt.ComputeReflect(input)

	if ccittComputeResult != ccittResult {
		t.Fatalf(`ccitt crc: %X doesnt match %X`, ccittComputeResult, ccittResult)
	}
}

// TODO init update result.
func TestAllSteps(t *testing.T) {
	//123456789123456789 in hex format
	input := []uint8{0x12, 0x34, 0x56, 0x78, 0x91, 0x23, 0x45, 0x67, 0x89}

	crcGsm := Crc16{
		Polynomial: 0x1021,
		StartMask:  0,
		EndMask:    0xFFFF,
	}
	gsmResult := uint16(0x64B3)

	crcGsm.Init()
	for _, v := range input {
		crcGsm.Update(v)
	}
	firstResult := crcGsm.Result()

	crcGsm.Init()
	for _, v := range input {
		crcGsm.Update(v)
	}
	secondResult := crcGsm.Result()

	if firstResult != gsmResult || secondResult != gsmResult {
		t.Fatalf(`gsm crc: %X doesnt match 1st: %X or 2nd: %X`, gsmResult, firstResult, secondResult)
	}

}
