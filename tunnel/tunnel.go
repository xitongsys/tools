package main

import (
	"crypto/cipher"
	"fmt"
	"net"
	"sync"
)

const (
	PACK_SIZE   = 1024
	BUFFER_SIZE = 4 * 1024
)

type Tunnel struct {
	// L: listener, C: client
	Role uint8
	Addr string

	TunConn     net.Conn
	InBuffer    []byte
	OutBuffer   []byte
	BufferMutex sync.Mutex

	ConnsCnt   uint64
	Conns      map[uint64]net.Conn
	ConnsMutex sync.Mutex

	CipherBlock cipher.Block

	Error error
}

func NewTunnel(role uint8, addr string, tunConn net.Conn, cipherBlock cipher.Block) *Tunnel {
	tun := &Tunnel{
		Role: role,
		Addr: addr,

		TunConn:   tunConn,
		InBuffer:  make([]byte, BUFFER_SIZE),
		OutBuffer: make([]byte, BUFFER_SIZE),

		ConnsCnt: 0,
		Conns:    map[uint64]net.Conn{},

		CipherBlock: cipherBlock,
	}

	return tun
}

func (tun *Tunnel) ReadMsg() (Msg, error) {
	return ReadMsg(tun.TunConn, tun.InBuffer, tun.CipherBlock)
}

func (tun *Tunnel) WriteMsg(msg Msg) error {
	tun.BufferMutex.Lock()
	defer tun.BufferMutex.Unlock()

	return WriteMsg(tun.TunConn, tun.OutBuffer, msg, tun.CipherBlock)
}

func (tun *Tunnel) NewId() uint64 {
	tun.ConnsMutex.Lock()
	defer tun.ConnsMutex.Unlock()

	tun.ConnsCnt++
	return tun.ConnsCnt
}

func (tun *Tunnel) Run() {
	if tun.Role == 'L' {
		tun.RunListener()
	} else {
		tun.RunClient()
	}
}

func (tun *Tunnel) RunListener() {
	go tun.TunHandler()

	listen, err := net.Listen("tcp", tun.Addr)
	if err != nil {
		return
	}
	defer listen.Close()

	msg := &MsgConn{}
	var conn net.Conn

	for tun.Error == nil {
		if conn, tun.Error = listen.Accept(); tun.Error == nil {
			msg.Id = tun.NewId()
			if tun.Error = tun.WriteMsg(msg); tun.Error == nil {
				tun.OpenConn(msg.Id, conn)
			}
		}
	}
}

func (tun *Tunnel) RunClient() {
	tun.TunHandler()
}

// open new connection
func (tun *Tunnel) OpenConn(id uint64, conn net.Conn) {
	tun.ConnsMutex.Lock()
	defer tun.ConnsMutex.Unlock()

	tun.Conns[id] = conn
	go tun.ConnHandler(conn, id)
}

// close connection
func (tun *Tunnel) CloseConn(id uint64, notify bool) {
	tun.ConnsMutex.Lock()
	conn, ok := tun.Conns[id]
	if ok {
		conn.Close()
		delete(tun.Conns, id)
	}
	tun.ConnsMutex.Unlock()

	if notify {
		msg := &MsgClose{
			Id: id,
		}
		tun.Error = tun.WriteMsg(msg)
	}
}

// conn -> tunnel
func (tun *Tunnel) ConnHandler(conn net.Conn, id uint64) {
	buf := make([]uint8, PACK_SIZE)
	msg := &MsgPack{
		Id: id,
	}

	var n int = 0
	var err error
	for tun.Error == nil && err == nil {
		if n, err = conn.Read(buf); n > 0 && err == nil {
			msg.DataLen = uint32(n)
			dataBufferLen := PaddlingLen(msg.DataLen)
			msg.Data = buf[:dataBufferLen]
			Paddling(buf[msg.DataLen:dataBufferLen])

			tun.Error = tun.WriteMsg(msg)

		} else if err != nil {
			tun.CloseConn(id, true)
		}
	}
}

// tunnel -> conn
func (tun *Tunnel) TunHandler() {
	var msgi Msg

	// listener
	if tun.Role == 'L' {
		for tun.Error == nil {
			if msgi, tun.Error = tun.ReadMsg(); tun.Error == nil {
				if msgi.Type() == PACK {
					msg := msgi.(*MsgPack)

					tun.ConnsMutex.Lock()
					conn, ok := tun.Conns[msg.Id]
					tun.ConnsMutex.Unlock()
					if ok {
						if _, err := conn.Write(msg.Data[:msg.DataLen]); err != nil {
							tun.CloseConn(msg.Id, true)
						}

					} else {
						tun.CloseConn(msg.Id, true)
					}

				} else if msgi.Type() == CLOSE {
					msg := msgi.(*MsgClose)
					tun.CloseConn(msg.Id, false)

				} else {
					tun.Error = fmt.Errorf("illegal msg type %v", msgi.Type())
				}
			}
		}

	} else { // client
		for tun.Error == nil {
			if msgi, tun.Error = tun.ReadMsg(); tun.Error == nil {
				if msgi.Type() == CONN {
					msg := msgi.(*MsgConn)
					if conn, err := net.Dial("tcp", tun.Addr); err == nil {
						tun.OpenConn(msg.Id, conn)

					} else {
						tun.CloseConn(msg.Id, true)
					}

				} else if msgi.Type() == PACK {
					msg := msgi.(*MsgPack)

					tun.ConnsMutex.Lock()
					conn, ok := tun.Conns[msg.Id]
					tun.ConnsMutex.Unlock()
					if ok {
						if _, err := conn.Write(msg.Data[:msg.DataLen]); err != nil {
							tun.CloseConn(msg.Id, true)
						}
					} else {
						tun.CloseConn(msg.Id, true)
					}

				} else if msgi.Type() == CLOSE {
					msg := msgi.(*MsgClose)
					tun.CloseConn(msg.Id, false)

				} else {
					tun.Error = fmt.Errorf("illegal msg type %v", msgi.Type())
				}
			}
		}
	}
}
