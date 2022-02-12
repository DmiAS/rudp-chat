package protocol

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// packet is a RUDP segment
type packet struct {
	seqNumber   uint32
	ackNumber   uint32
	segmentType rudpSegmentType
	payload     []byte
}

type rudpSegmentType int8

const (
	rudpSegmentTypeNormal rudpSegmentType = iota
	rudpSegmentTypeConn
	rudpSegmentTypeConnAck
	rudpSegmentTypeFin
	rudpSegmentTypeAck
	rudpSegmentTypePin
	rudpSegmentTypeErr
)

// ------------------------------------------
// |          Seq Number(32bits)            |
// ------------------------------------------
// |          Ack Number(32bits)            |
// ------------------------------------------
// |C|A|F|P|E|12bits...|.......2byte........|
// |O|C|I|I|R|reserved.|.....data len.......|
// |N|K|N|N|R|.........|....................|
// -----------------------------------------
// |................data....................|

func unmarshalRUDPPacket(data []byte) (*packet, error) {
	r := &packet{}
	if len(data) < rudpHeaderLen {
		return nil, errors.New("illegal rudp segment")
	}

	r.seqNumber = binary.BigEndian.Uint32(data[0:4])
	r.ackNumber = binary.BigEndian.Uint32(data[4:8])

	connFlag, ackFlag, finFlag, pinFlag, errFlag := parseFlag(data[8])
	switch {
	case connFlag:
		if ackFlag {
			r.segmentType = rudpSegmentTypeConnAck
		} else {
			r.segmentType = rudpSegmentTypeConn
		}
	case finFlag:
		r.segmentType = rudpSegmentTypeFin
	case pinFlag:
		r.segmentType = rudpSegmentTypePin
	case ackFlag:
		r.segmentType = rudpSegmentTypeAck
	case errFlag:
		r.segmentType = rudpSegmentTypeErr
	default:
		r.segmentType = rudpSegmentTypeNormal
	}

	payloadLen := (binary.BigEndian.Uint16(data[10:12]))
	r.payload = data[rudpHeaderLen : rudpHeaderLen+int(payloadLen)]
	return r, nil
}

func parseFlag(flag byte) (connFlag, ackFlag, finFlag, pinFlag, errFlag bool) {
	if flag&(1<<7) != 0 {
		connFlag = true
	} else {
		connFlag = false
	}

	if flag&(1<<6) != 0 {
		ackFlag = true
	} else {
		ackFlag = false
	}

	if flag&(1<<5) != 0 {
		finFlag = true
	} else {
		finFlag = false
	}

	if flag&(1<<4) != 0 {
		pinFlag = true
	} else {
		pinFlag = false
	}

	if flag&(1<<3) != 0 {
		errFlag = true
	} else {
		errFlag = false
	}
	return
}

func (p *packet) marshal() []byte {
	buf := make([]byte, rudpHeaderLen+len(p.payload))
	binary.BigEndian.PutUint32(buf[0:4], p.seqNumber)
	binary.BigEndian.PutUint32(buf[4:8], p.ackNumber)
	var flag byte
	switch p.segmentType {
	case rudpSegmentTypeNormal:
		flag = 0
	case rudpSegmentTypeConn:
		flag |= (1 << 7)
	case rudpSegmentTypeConnAck:
		flag |= (1 << 7)
		flag |= (1 << 6)
	case rudpSegmentTypeAck:
		flag |= (1 << 6)
	case rudpSegmentTypeFin:
		flag |= (1 << 5)
	case rudpSegmentTypePin:
		flag |= (1 << 4)
	case rudpSegmentTypeErr:
		flag |= (1 << 3)
	}
	buf[8] = flag
	binary.BigEndian.PutUint16(buf[10:12], uint16(len(p.payload)))
	for i, b := range p.payload {
		buf[rudpHeaderLen+i] = b
	}
	return buf
}

func newNormalPacket(payload []byte, seqNumber uint32) *packet {
	return &packet{
		seqNumber:   seqNumber,
		segmentType: rudpSegmentTypeNormal,
		payload:     payload,
	}
}

func newAckPacket(ackNumber uint32) *packet {
	return &packet{
		ackNumber:   ackNumber,
		segmentType: rudpSegmentTypeAck,
	}
}

func newConPacket(seqNumber uint32) *packet {
	return &packet{
		seqNumber:   seqNumber,
		segmentType: rudpSegmentTypeConn,
	}
}

func newConAckPacket(ackNumber uint32) *packet {
	return &packet{
		seqNumber:   ackNumber,
		segmentType: rudpSegmentTypeConnAck,
	}
}

func newFinPacket() *packet {
	return &packet{
		segmentType: rudpSegmentTypeFin,
	}
}

func newPinPacket() *packet {
	return &packet{
		segmentType: rudpSegmentTypePin,
	}
}

func (p *packet) print() {
	if debug {
		fmt.Printf(
			"packet: seqNumber = %d#ackNumber = %d#segmentType = %d#payload = %s\n", p.seqNumber, p.ackNumber,
			p.segmentType, string(p.payload),
		)
	}
}
