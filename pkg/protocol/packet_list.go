package protocol

import (
	"fmt"
	"sync"
	"time"
)

type packetListOrderType int8

const (
	packetListOrderBySeqNb packetListOrderType = iota
	packetListOrderByAckNb
)

// packetList is a sorted link list of packet that ordered by packet.SeqNumber/AckNumber ASC
type packetList struct {
	head, rail *node
	length     int32
	orderType  packetListOrderType
	cond       *sync.Cond
	mutex      *sync.Mutex
}

type node struct {
	data   *packet
	next   *node
	tsNano int64
}

func newPacketList(orderType packetListOrderType) *packetList {
	r := &packetList{
		orderType: orderType,
		mutex:     &sync.Mutex{},
	}
	r.cond = sync.NewCond(r.mutex)
	return r
}

func (l *packetList) getPacketSortKey(p *packet) uint32 {
	if l.orderType == packetListOrderBySeqNb {
		return p.seqNumber
	} else { // if l.orderType == packetListOrderByAckNb
		return p.ackNumber
	}
}

func (l *packetList) putPacket(p *packet) {
	newNode := &node{data: p, tsNano: time.Now().UnixNano()}
	l.mutex.Lock()
	defer func() {
		l.cond.Signal()
		l.mutex.Unlock()
	}()
	if l.head == nil {
		l.head = newNode
		l.rail = newNode
		l.length++
	} else {
		var last *node
		cur := l.head
		sortKey := l.getPacketSortKey(p)
		for ; cur != nil; cur = cur.next {
			curSortKey := l.getPacketSortKey(cur.data)
			if curSortKey == sortKey {
				return
			} else if curSortKey > sortKey {
				newNode.next = cur
				if last == nil { // only one node now
					l.head = newNode
					l.rail = cur
				} else {
					last.next = newNode
				}
				l.length++
				return
			} else {
				last = cur
			}
		}
		l.rail.next = newNode
		l.rail = newNode
		l.length++
	}
}

func (l *packetList) empty() bool {
	return l.length == 0
}

func (l *packetList) consume() *packet {
	l.mutex.Lock()
	defer func() {
		l.length--
		l.mutex.Unlock()
	}()

	for {
		if l.head != nil {
			head := l.head
			if l.head == l.rail {
				l.head = nil
				l.rail = nil
			} else {
				l.head = head.next
			}
			return head.data
		} else {
			l.cond.Wait()
		}
	}
}

func (l *packetList) consumePacketSinceNMs(N int) []*packet {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.empty() {
		return nil
	}
	// res should be nil at most time
	var res []*packet
	now := time.Now().UnixNano()
	cur := l.head
	for ; cur != nil; cur = cur.next {
		dis := now - cur.tsNano
		if dis >= int64(N*1e6) {
			res = append(res, cur.data)
		} else {
			break
		}
	}
	l.head = cur
	if cur == nil {
		l.rail = nil
	}
	l.length -= int32(len(res))
	return res
}

// removePacketBySeqNb return val means if the seqNb packet is be found and deleted
func (l *packetList) removePacketByNb(nb uint32) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.empty() {
		return
	}
	var last *node
	for cur := l.head; cur != nil; cur = cur.next {
		curNb := l.getPacketSortKey(cur.data)
		if curNb == nb {
			l.length--
			if cur == l.head {
				if l.head == l.rail {
					l.head = nil
					l.rail = nil
				} else {
					l.head = l.head.next
				}
			} else {
				last.next = cur.next
				if cur == l.rail {
					l.rail = last
				}
			}
			return
		}
		last = cur
	}
}

func (l *packetList) debug() {
	if l.head == nil {
		fmt.Println("head is nil!!")
		return
	}
	output := "############################\n"
	for cur := l.head; cur != nil; cur = cur.next {
		s := fmt.Sprintf("packetInTheList: seqNb: %d, payload: %s\n", cur.data.seqNumber, string(cur.data.payload))
		output += s
	}
	output += "############################\n"
	fmt.Print(output)
}
