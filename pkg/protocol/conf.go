package protocol

import "fmt"

var (
	rawUDPPacketLenLimit    = 66528 // default is the maximum for raw UDP
	rudpHeaderLen           = 12
	rudpPayloadLenLimit     = rawUDPPacketLenLimit - rudpHeaderLen
	resendDelayThreshholdMS = 5 // 3 ms
)

const (
	defaultSendTickNano         = 1e7 // 10 ms
	defaultHeartBeatCycleMinute = 30
)

const (
	RawUDPSendNotComplete = "raw udp not send the complete rudp packet"
)

var debug = false

func Debug() {
	debug = true
}

func log(format string, a ...interface{}) {
	if debug {
		fmt.Printf(format, a...)
	}
}

func SetRawUDPPacketLenLimit(size int) {
	rawUDPPacketLenLimit = size
}
