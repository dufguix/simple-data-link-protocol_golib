package linklayerprotocol

type AppReceiver interface {
	Receive(payload []uint8)
}

type ComSender interface {
	Send(bytes []uint8)
}
