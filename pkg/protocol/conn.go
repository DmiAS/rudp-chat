package protocol

import (
	"errors"
	"io"
	"math/rand"
	"net"
	"time"
)

const (
	maxWaitSegmentCntWhileConn = 3
)

var rander = rand.New(rand.NewSource(time.Now().Unix()))

type connStatus int8

const (
	connStatusConnecting connStatus = iota
	connStatusOpen
	connStatusClose
	connStatusErr
)

type connType int8

const (
	connTypeServer connType = iota
	connTypeClient
)

var (
	rudpConnClosedErr     = errors.New("the rudp connection closed")
	resolveRUDPSegmentErr = errors.New("rudp message segment resolved error")
)

// RUDPConn is the reliable conn base on udp
type RUDPConn struct {
	sendPacketChannel chan *packet // all kinds of packet already to send
	resendPacketQueue *packetList  // packet from sendPacketQueue that waiting for the peer's ack segment

	recvPacketChannel   chan *packet // all kinds of packets recv: ack/conn/fin/etc....
	outputPacketQueue   *packetList  // the data for application layer
	outputDataTmpBuffer []byte       // temporary save the surplus data for aplication layer that would be read next time

	rawUDPDataChan chan []byte // only server conn need to

	localAddr  *net.UDPAddr
	remoteAddr *net.UDPAddr
	rawUDPConn *net.UDPConn

	sendSeqNumber       uint32
	lastRecvTs          int64 // last recv data unix timestamp
	lastSendTs          int64
	maxHasReadSeqNumber uint32 // the max number that the application has read

	sendTickNano int32

	sendTickModifyEvent chan int32

	heartBeatCycleMinute int8

	closeConnCallback         func() // execute when connection closed
	buildConnCallbackListener func() // execute when connection build
	closeConnCallbackListener func()

	rudpConnStatus connStatus
	rudpConnType   connType

	// react by errBus
	sendStop          chan error
	resendStop        chan error
	recvStop          chan error
	packetHandlerStop chan error
	heartbeatStop     chan error

	errBus chan error

	err error
}

// DialRUDP client dial server, building a relieable connection
func DialRUDP(conn *net.UDPConn, remoteAddr *net.UDPAddr) (*RUDPConn, error) {
	c := &RUDPConn{}
	c.localAddr = nil
	c.remoteAddr = remoteAddr
	c.rudpConnType = connTypeClient
	c.sendSeqNumber = 0
	c.recvPacketChannel = make(chan *packet, 1<<5)
	c.rudpConnStatus = connStatusConnecting

	if err := c.clientBuildConn(conn); err != nil {
		return nil, err
	}

	c.sendPacketChannel = make(chan *packet, 1<<5)
	c.resendPacketQueue = newPacketList(packetListOrderBySeqNb)
	c.outputPacketQueue = newPacketList(packetListOrderBySeqNb)

	c.rudpConnStatus = connStatusOpen
	c.sendTickNano = defaultSendTickNano
	c.sendTickModifyEvent = make(chan int32, 1)
	c.heartBeatCycleMinute = defaultHeartBeatCycleMinute
	c.lastSendTs = time.Now().Unix()
	c.lastRecvTs = time.Now().Unix()

	c.sendStop = make(chan error, 1)
	c.recvStop = make(chan error, 1)
	c.resendStop = make(chan error, 1)
	c.packetHandlerStop = make(chan error, 1)
	c.heartbeatStop = make(chan error, 1)
	c.errBus = make(chan error, 1)

	c.closeConnCallback = func() {
		c.rudpConnStatus = connStatusClose
		c.errBus <- io.EOF
	}

	// monitor errBus
	go c.errWatcher()

	// net io
	go c.recv()
	go c.send()

	go c.resend()
	go c.packetHandler()

	// client need to keep a heart beat
	go c.keepLive()
	log("build the RUDP connection succ!\n")
	return c, nil
}

func (c *RUDPConn) errWatcher() {
	err := <-c.errBus
	log("errBus recv error: %v\n", err)
	c.err = err
	c.resendStop <- err
	c.recvStop <- err
	c.sendStop <- err
	if c.rudpConnType == connTypeClient {
		c.heartbeatStop <- err
	}

	if c.rudpConnType == connTypeServer {
		c.closeConnCallbackListener()
	}

}

func (c *RUDPConn) keepLive() {
	ticker := time.NewTicker(time.Duration(c.heartBeatCycleMinute) * 60 * time.Second)
	defer ticker.Stop()
	select {
	case <-ticker.C:
		now := time.Now().Unix()
		if now-c.lastSendTs >= int64(c.heartBeatCycleMinute)*60 {
			c.sendPacketChannel <- newPinPacket()
		}
	case <-c.heartbeatStop:
		return
	}
}

func (c *RUDPConn) recv() {
	if c.rudpConnType == connTypeClient {
		c.clientRecv()
	} else if c.rudpConnType == connTypeServer {
		c.serverRecv()
	}
}

func (c *RUDPConn) clientRecv() {
	for {
		select {
		case <-c.recvStop:
			return
		default:
			buf := make([]byte, rawUDPPacketLenLimit)
			n, err := c.rawUDPConn.Read(buf)
			if err != nil {
				c.errBus <- err
				return
			}
			buf = buf[:n]
			apacket, err := unmarshalRUDPPacket(buf)
			if err != nil {
				c.errBus <- resolveRUDPSegmentErr
				return
			}
			c.recvPacketChannel <- apacket
		}
	}
}

func (c *RUDPConn) serverRecv() {
	for {
		select {
		case <-c.recvStop:
			return
		case data := <-c.rawUDPDataChan:
			apacket, err := unmarshalRUDPPacket(data)
			if err != nil {
				c.errBus <- resolveRUDPSegmentErr
				return
			}
			c.recvPacketChannel <- apacket
		}
	}
}

// handle the recv packets
func (c *RUDPConn) packetHandler() {
	for {
		select {
		case <-c.packetHandlerStop:
			return
		case apacket := <-c.recvPacketChannel:
			switch apacket.segmentType {
			case rudpSegmentTypeNormal:
				if apacket.seqNumber <= c.maxHasReadSeqNumber {
					// discard
					continue
				}
				c.outputPacketQueue.putPacket(apacket)
				// ack
				c.sendPacketChannel <- newAckPacket(apacket.seqNumber)
			case rudpSegmentTypeAck:
				log("ack %d\n", apacket.ackNumber)
				c.resendPacketQueue.removePacketByNb(apacket.ackNumber)
			case rudpSegmentTypeFin:
				c.errBus <- io.EOF
				return
			case rudpSegmentTypePin:
				// do nothing
			case rudpSegmentTypeConn:
				if c.rudpConnType != connTypeServer {
					continue
				}
				// server send CON ack segment
				segment := newConAckPacket(apacket.seqNumber).marshal()
				n, err := c.write(segment)
				if err != nil {
					c.errBus <- err
					return
				}
				if n != len(segment) {
					c.errBus <- errors.New(RawUDPSendNotComplete)
					return
				}
				// build conn
				log("server send CON-ACK segment\n")
				c.rudpConnStatus = connStatusOpen
				c.buildConnCallbackListener()
			}
			c.lastRecvTs = time.Now().Unix()
		}
	}
}

func (c *RUDPConn) send() {
	ticker := time.NewTicker(time.Duration(c.sendTickNano) * time.Nanosecond)
	defer ticker.Stop()
	for {
		select {
		case c.sendTickNano = <-c.sendTickModifyEvent:
			ticker.Stop()
			ticker = time.NewTicker(time.Duration(c.sendTickNano) * time.Nanosecond)
		case <-c.sendStop:
			return
		case <-ticker.C:
			c.sendPacket()
			c.lastSendTs = time.Now().Unix()
		}
	}
}

// SetRealSendTick modify the segment sending cycle
func (c *RUDPConn) SetSendTick(nano int32) {
	c.sendTickModifyEvent <- nano
}

func (c *RUDPConn) write(data []byte) (n int, err error) {
	n, err = c.rawUDPConn.WriteTo(data, c.remoteAddr)
	return
}

func (c *RUDPConn) sendPacket() {
	apacket := <-c.sendPacketChannel
	segment := apacket.marshal()
	n, err := c.write(segment)
	if err != nil {
		log("sendPacket error: %v, %d", err, len(segment))
		c.errBus <- err
		return
	}
	if n != len(segment) {
		c.errBus <- errors.New(RawUDPSendNotComplete)
		return
	}
	// apacket.print()
	// only the normal segment possiblely needs to resend
	if apacket.segmentType == rudpSegmentTypeNormal {
		c.resendPacketQueue.putPacket(apacket)
	}
}

func (c *RUDPConn) clientBuildConn(udpConn *net.UDPConn) error {
	// just init instance
	c.rawUDPConn = udpConn
	c.rawUDPConn.SetWriteBuffer(65528)
	// send conn segment
	connSeqNb := c.sendSeqNumber
	c.sendSeqNumber++
	connSegment := newConPacket(connSeqNb).marshal()
	n, err := udpConn.WriteTo(connSegment, c.remoteAddr)
	if err != nil {
		return err
	}
	if n != len(connSegment) {
		return errors.New(RawUDPSendNotComplete)
	}
	log("client send the CONN segment\n")

	// wait the server ack conn segment
	// may the server's ack segment and normal segment out-of-order
	// so if the recv not the ack segment, we try to wait the next
	for cnt := 0; cnt < maxWaitSegmentCntWhileConn; cnt++ {
		buf := make([]byte, rawUDPPacketLenLimit)
		n, err = udpConn.Read(buf)
		if err != nil {
			return err
		}
		recvPacket, err := unmarshalRUDPPacket(buf[:n])
		if err != nil {
			return errors.New("analyze the recvSegment error: " + err.Error())
		}

		if recvPacket.ackNumber == connSeqNb && recvPacket.segmentType == rudpSegmentTypeConnAck {
			// conn OK
			log("client recv the server CON-ACK segment\n")
			return nil
		} else {
			// c.recvPacketChannel <- recvPacket
			continue
		}
	}
	return nil
}

func serverBuildConn(rawUDPConn *net.UDPConn, remoteAddr *net.UDPAddr) (*RUDPConn, error) {
	c := &RUDPConn{}
	c.rawUDPConn = rawUDPConn
	c.rawUDPConn.SetWriteBuffer(65528)
	c.localAddr, _ = net.ResolveUDPAddr(rawUDPConn.LocalAddr().Network(), rawUDPConn.LocalAddr().String())
	c.remoteAddr = remoteAddr
	c.rudpConnType = connTypeServer
	c.sendSeqNumber = 0
	c.rudpConnStatus = connStatusConnecting

	c.recvPacketChannel = make(chan *packet, 1<<5)
	c.sendPacketChannel = make(chan *packet, 1<<5)
	c.rawUDPDataChan = make(chan []byte, 1<<5)
	c.resendPacketQueue = newPacketList(packetListOrderBySeqNb)
	c.outputPacketQueue = newPacketList(packetListOrderBySeqNb)

	c.sendTickNano = defaultSendTickNano
	c.sendTickModifyEvent = make(chan int32, 1)
	c.lastRecvTs = time.Now().Unix()

	c.sendStop = make(chan error, 1)
	c.recvStop = make(chan error, 1)
	c.resendStop = make(chan error, 1)
	c.packetHandlerStop = make(chan error, 1)
	c.errBus = make(chan error, 1)

	c.closeConnCallback = func() {
		c.rudpConnStatus = connStatusClose
		c.errBus <- io.EOF
	}

	go c.errWatcher()

	// net io
	go c.send()
	go c.recv()

	go c.resend()
	go c.packetHandler()

	return c, nil
}

func (c *RUDPConn) Read(b []byte) (int, error) {
	readCnt := len(b)
	n := len(b)
	if n == 0 {
		return 0, nil
	}
	curWrite := 0
	if len(c.outputDataTmpBuffer) != 0 {
		if n <= len(c.outputDataTmpBuffer) {
			copy(b, c.outputDataTmpBuffer[:n])
			c.outputDataTmpBuffer = c.outputDataTmpBuffer[n:]
			return readCnt, nil
		} else {
			n -= len(c.outputDataTmpBuffer)
			curWrite += len(c.outputDataTmpBuffer)
			copy(b, c.outputDataTmpBuffer)
		}
	}

	for n > 0 {
		apacket := c.outputPacketQueue.consume()
		if apacket.seqNumber-c.maxHasReadSeqNumber != 1 {
			log("发生丢包 cur %d max %d\n", apacket.seqNumber, c.maxHasReadSeqNumber)
		}
		// apacket.print()
		c.maxHasReadSeqNumber = apacket.seqNumber
		data := apacket.payload
		if n <= len(data) {
			copy(b[curWrite:], data[:n])
			c.outputDataTmpBuffer = data[n:]
			return readCnt, nil
		} else {
			copy(b[curWrite:], data)
			n -= len(data)
			curWrite += len(data)
		}
	}
	return 0, nil
}

func (c *RUDPConn) Write(b []byte) (int, error) {
	if c.err != nil {
		return 0, errors.New("rudp write error: " + c.err.Error())
	}
	n := len(b)
	for {
		if len(b) <= rudpPayloadLenLimit {
			c.sendPacketChannel <- newNormalPacket(b, c.sendSeqNumber)
			c.sendSeqNumber++
			return n, nil
		} else {
			c.sendPacketChannel <- newNormalPacket(b[:rudpPayloadLenLimit], c.sendSeqNumber)
			c.sendSeqNumber++
			b = b[rudpPayloadLenLimit:]
		}
	}
}

// Close close must be called while closing the conn
func (c *RUDPConn) Close() error {
	if c.rudpConnStatus != connStatusOpen {
		return errors.New("the rudp connection is not open status!")
	}
	defer func() {
		if c.rudpConnType == connTypeServer {
			c.closeConnCallbackListener()
		}
		c.closeConnCallback()
	}()

	finSegment := newFinPacket().marshal()
	n, err := c.write(finSegment)
	if err != nil {
		return err
	}
	if n != len(finSegment) {
		return errors.New(RawUDPSendNotComplete)
	}
	c.errBus <- io.EOF
	return nil
}

func (c *RUDPConn) resend() {
	ticker := time.NewTicker(time.Millisecond * time.Duration(resendDelayThreshholdMS))
	defer ticker.Stop()
	for {
		select {
		case <-c.resendStop:
			return
		case <-ticker.C:
			resendPacketList := c.resendPacketQueue.consumePacketSinceNMs(resendDelayThreshholdMS)
			if len(resendPacketList) != 0 {
				log("一轮重传\n")
			}
			for _, resendPacket := range resendPacketList {
				segment := resendPacket.marshal()
				n, err := c.write(segment)
				if err != nil {
					c.errBus <- err
					return
				}
				if n != len(segment) {
					c.errBus <- errors.New(RawUDPSendNotComplete)
					return
				}
				log("重传 %d\n", resendPacket.seqNumber)
			}
		}
	}
}

func (c *RUDPConn) LocalAddr() net.Addr {
	return c.localAddr
}

func (c *RUDPConn) RemoteAddr() net.Addr {
	return c.remoteAddr
}

func (c *RUDPConn) SetDeadline(t time.Time) error {
	return nil
}

func (c *RUDPConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *RUDPConn) SetWriteDeadline(t time.Time) error {
	return nil
}
