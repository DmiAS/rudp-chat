package rudp

const (
	// key and connection parameters
	pass         = "pass"
	salt         = "salt"
	iter         = 1024
	keyLen       = 32
	dataShards   = 10
	parityShards = 3

	// chan sizes
	recvChanSize = 5
	sendChanSize = 5

	// buffer size for both type os messages
	bufferSize = 64000
)
