package protocol

import (
	"errors"
	"net"
	"sync"
)

// RUDPListener implements net.Listener interface
type RUDPListener struct {
	listenConn   *net.UDPConn
	peerMap      *sync.Map
	newConnQueue chan *RUDPConn
	err          chan error
}

func ListenRUDP(conn *net.UDPConn) *RUDPListener {
	listener := &RUDPListener{
		listenConn:   conn,
		peerMap:      &sync.Map{},
		err:          make(chan error, 1),
		newConnQueue: make(chan *RUDPConn, 1<<10),
	}
	go listener.listen()
	return listener
}

func (l *RUDPListener) Accept() (net.Conn, error) {
	select {
	case c, ok := <-l.newConnQueue:
		if !ok {
			return nil, errors.New("listener closed")
		}
		return c, nil
	case err := <-l.err:
		return nil, err
	}
}

func (l *RUDPListener) Close() error {

	return nil
}

func (l *RUDPListener) Addr() net.Addr {
	return l.listenConn.LocalAddr()
}

func (l *RUDPListener) listen() {
	log("listen on %s\n", l.Addr().String())
	for {
		buf := make([]byte, rawUDPPacketLenLimit)
		n, remoteAddr, err := l.listenConn.ReadFromUDP(buf)
		if err != nil {
			l.err <- err
			return
		}
		buf = buf[:n]
		if v, ok := l.peerMap.Load(remoteAddr.String()); ok {
			rudpConn := v.(*RUDPConn)
			rudpConn.rawUDPDataChan <- buf
		} else {
			rudpConn, err := serverBuildConn(l.listenConn, remoteAddr)
			if err != nil {
				l.err <- err
				return
			}
			rudpConn.buildConnCallbackListener = func() {
				log("accept new RUDP connection from %v\n", rudpConn.RemoteAddr())
				l.newConnQueue <- rudpConn
			}
			rudpConn.closeConnCallbackListener = func() {
				l.peerMap.Delete(rudpConn.remoteAddr)
			}
			l.peerMap.Store(remoteAddr.String(), rudpConn)
			rudpConn.rawUDPDataChan <- buf
		}
	}
}
