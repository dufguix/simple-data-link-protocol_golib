package linklayerprotocol

import (
	"fmt"
	"testing"
)

type FakeAppReceiver struct {
	callback func([]uint8)
}

func (far *FakeAppReceiver) Receive(payload []uint8) {
	far.callback(payload)
}

type FakeComSender struct {
	callback func([]uint8)
}

func (fcs *FakeComSender) Send(payload []uint8) {
	fcs.callback(payload)
}

func testEq(a, b []uint8) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestGoodReceiving(t *testing.T) {
	//setup
	crcXmodem := Crc16{
		Polynomial: 0x1021,
		StartMask:  0,
		EndMask:    0,
	}
	xmodemResult := uint16(0x9B4C)
	const length uint8 = 0x09
	input := []uint8{StartByte, length, 0x12, 0x34, 0x56, 0x78, 0x91, 0x23, 0x45, 0x67, 0x89, uint8(xmodemResult >> 8), uint8(xmodemResult)}
	output := []uint8{0x12, 0x34, 0x56, 0x78, 0x91, 0x23, 0x45, 0x67, 0x89}
	var outputResult []uint8

	callbackCalled := 0
	fakeReceiver := FakeAppReceiver{callback: func(p []uint8) {
		callbackCalled++
		outputResult = p
	}}

	linkP := NewDefaultLinkProtocol(nil, &fakeReceiver, crcXmodem)

	// test begin here
	for _, v := range input {
		linkP.ReceiveByte(v)
	}

	if callbackCalled != 1 {
		t.Errorf("AppReceiver.receive() not called")
	}

	if !testEq(outputResult, output) {
		for _, v := range outputResult {
			fmt.Printf("%X", v)
		}
		print("\n")
		for _, v := range output {
			fmt.Printf("%X", v)
		}
		print("\n")
		t.Errorf("payload is not correctly extracted")

	}

	// run twice for testing internal reset
	for _, v := range input {
		linkP.ReceiveByte(v)
	}

	if callbackCalled != 2 {
		t.Errorf("AppReceiver.receive() not called")
	}

	if !testEq(outputResult, output) {
		for _, v := range outputResult {
			fmt.Printf("%X", v)
		}
		print("\n")
		for _, v := range output {
			fmt.Printf("%X", v)
		}
		print("\n")
		t.Errorf("payload is not correctly extracted")
	}

}

func TestBadReceiving(t *testing.T) {
	//setup
	crcXmodem := Crc16{
		Polynomial: 0x1021,
		StartMask:  0,
		EndMask:    0,
	}
	xmodemResult := uint16(0x9B4C)
	const length uint8 = 0x09
	//last byte 0x89 becomes 0x88
	input := []uint8{StartByte, length, 0x12, 0x34, 0x56, 0x78, 0x91, 0x23, 0x45, 0x67, 0x88, uint8(xmodemResult >> 8), uint8(xmodemResult)}
	//var outputResult []uint8

	callbackCalled := 0
	fakeReceiver := FakeAppReceiver{callback: func(p []uint8) {
		callbackCalled++
		//outputResult = p
	}}

	linkP := NewDefaultLinkProtocol(nil, &fakeReceiver, crcXmodem)

	// test begin here
	for _, v := range input {
		linkP.ReceiveByte(v)
	}

	if callbackCalled != 0 {
		t.Errorf("AppReceiver.receive() should not be called")
	}
}

func TestSend(t *testing.T) {
	//setup
	crcXmodem := Crc16{
		Polynomial: 0x1021,
		StartMask:  0,
		EndMask:    0,
	}
	xmodemResult := uint16(0x9B4C)
	const length uint8 = 0x09
	//last byte 0x89 becomes 0x88
	input := []uint8{0x12, 0x34, 0x56, 0x78, 0x91, 0x23, 0x45, 0x67, 0x89}
	output := []uint8{StartByte, length, 0x12, 0x34, 0x56, 0x78, 0x91, 0x23, 0x45, 0x67, 0x89, uint8(xmodemResult >> 8), uint8(xmodemResult)}
	var outputResult []uint8

	callbackCalled := 0
	fakeComSender := FakeComSender{callback: func(b []uint8) {
		callbackCalled++
		outputResult = b
	}}

	linkP := NewDefaultLinkProtocol(&fakeComSender, nil, crcXmodem)

	// test begin here

	linkP.Send(input)
	for _, v := range input {
		linkP.ReceiveByte(v)
	}

	if callbackCalled != 1 {
		t.Errorf("ComSender.send() not called")
	}

	if !testEq(outputResult, output) {
		println(outputResult)
		for _, v := range outputResult {
			fmt.Printf("%X", v)
		}
		print("\n")
		println(output)
		for _, v := range output {
			fmt.Printf("%X", v)
		}
		print("\n")
		t.Errorf("frame is not correctly built")
	}
}
