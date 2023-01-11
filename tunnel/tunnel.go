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

type Listen struct {
	Id          uint64
	Listener    net.Listener
	ListenAddr  string
	ForwardAddr string
}

func (listen *Listen) String() string {
	return fmt.Sprintf("{Id:%v, ListenAddr:%v, ForwardAddr:%v}", listen.Id, listen.ListenAddr, listen.ForwardAddr)
}

type Connection struct {
	Id   uint64
	Conn net.Conn
}

func (conn *Connection) String() string {
	return fmt.Sprintf("{Id:%v, RemoteAddr:%v}", conn.Id, conn.Conn.RemoteAddr().String())
}

type Tunnel struct {
	Name       string
	RemoteAddr string
	Password   string

	TunConn     net.Conn
	InBuffer    []byte
	OutBuffer   []byte
	BufferMutex sync.Mutex

	Listens      map[uint64]*Listen
	ListensMutex sync.Mutex

	Conns      map[uint64]*Connection
	ConnsMutex sync.Mutex

	IdCnt      uint64
	IdCntMutex sync.Mutex

	CipherBlock cipher.Block

	Error error
}

func NewTunnel(name string, remoteAddr string, password string, tunConn net.Conn) *Tunnel {
	tun := &Tunnel{
		Name:       name,
		RemoteAddr: remoteAddr,
		Password:   password,
		TunConn:    tunConn,
		InBuffer:   make([]byte, BUFFER_SIZE),
		OutBuffer:  make([]byte, BUFFER_SIZE),

		Listens: map[uint64]*Listen{},
		Conns:   map[uint64]*Connection{},

		IdCnt: 0,

		CipherBlock: Password2Cipher(password),
	}

	return tun
}

func (tun *Tunnel) Exit(err error) {
	tun.TunConn.Close()
	tun.Error = err
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
	tun.IdCntMutex.Lock()
	defer tun.IdCntMutex.Unlock()

	tun.IdCnt++
	return tun.IdCnt
}

// open new listen
func (tun *Tunnel) OpenListen(id uint64, listenAddr, forwardAddr string) {
	tun.ListensMutex.Lock()
	defer tun.ListensMutex.Unlock()

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return
	}

	tun.Listens[id] = &Listen{
		Id:          id,
		Listener:    listener,
		ListenAddr:  listenAddr,
		ForwardAddr: forwardAddr,
	}

	forwardAddrInt, err := Addr2Int(forwardAddr)
	if err != nil {
		return
	}

	go func() {
		for tun.Error == nil {
			if conn, err := listener.Accept(); err == nil {
				msg := &MsgConn{
					Id:   tun.NewId(),
					Addr: forwardAddrInt,
				}

				if tun.Error = tun.WriteMsg(msg); tun.Error == nil {
					tun.OpenConn(msg.Id, conn)
				}

			} else {
				break
			}
		}
	}()
}

// close listen
func (tun *Tunnel) CloseListen(id uint64) {
	tun.ListensMutex.Lock()
	defer tun.ListensMutex.Unlock()
	tun.closeListen(id)
}

// close listen without lock
func (tun *Tunnel) closeListen(id uint64) {
	if listen, ok := tun.Listens[id]; ok {
		listen.Listener.Close()
	}
}

// open new connection
func (tun *Tunnel) OpenConn(id uint64, conn net.Conn) {
	tun.ConnsMutex.Lock()
	defer tun.ConnsMutex.Unlock()

	tun.Conns[id] = &Connection{
		Id:   id,
		Conn: conn,
	}
	go tun.ConnHandler(conn, id)
}

// close connection with out lock
func (tun *Tunnel) closeConn(id uint64, notify bool) {
	conn, ok := tun.Conns[id]
	if ok {
		conn.Conn.Close()
		delete(tun.Conns, id)
	}

	if notify {
		msg := &MsgCloseConn{
			Id: id,
		}
		if err := tun.WriteMsg(msg); err != nil {
			tun.Error = err
		}
	}
}

// close connection
func (tun *Tunnel) CloseConn(id uint64, notify bool) {
	tun.ConnsMutex.Lock()
	defer tun.ConnsMutex.Unlock()
	tun.closeConn(id, notify)
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

			if err := tun.WriteMsg(msg); err != nil {
				tun.Error = err
			}

		} else if err != nil {
			tun.CloseConn(id, true)
		}
	}
}

// tunnel -> conn
func (tun *Tunnel) Run() {
	var msgi Msg
	var err error

	for tun.Error == nil {
		if msgi, err = tun.ReadMsg(); err == nil {
			if msgi.Type() == PACK {
				msg := msgi.(*MsgPack)

				tun.ConnsMutex.Lock()
				conn, ok := tun.Conns[msg.Id]
				tun.ConnsMutex.Unlock()
				if ok {
					if _, err := conn.Conn.Write(msg.Data[:msg.DataLen]); err != nil {
						tun.CloseConn(msg.Id, true)
					}

				} else {
					tun.CloseConn(msg.Id, true)
				}

			} else if msgi.Type() == LISTEN {
				msg := msgi.(*MsgListen)

				listenAddr, err1 := Int2Addr(msg.ListenAddr)
				forwardAddr, err2 := Int2Addr(msg.ForwardAddr)

				if err1 == nil && err2 == nil {
					tun.OpenListen(tun.NewId(), listenAddr, forwardAddr)
				}

			} else if msgi.Type() == CLOSELISTEN {
				msg := msgi.(*MsgCloseConn)
				tun.CloseListen(msg.Id)

			} else if msgi.Type() == CLOSECONN {
				msg := msgi.(*MsgCloseConn)
				tun.CloseConn(msg.Id, false)

			} else if msgi.Type() == CONN {
				msg := msgi.(*MsgConn)
				addr, err := Int2Addr(msg.Addr)

				if err != nil {
					tun.CloseConn(msg.Id, true)

				} else {
					if conn, err := net.Dial("tcp", addr); err == nil {
						tun.OpenConn(msg.Id, conn)

					} else {
						tun.CloseConn(msg.Id, true)
					}
				}

			} else {
				tun.Error = fmt.Errorf("illegal msg type %v", msgi.Type())
			}

		} else {
			tun.Error = err
		}
	}
}

func (tun *Tunnel) CloseListens() {
	tun.ListensMutex.Lock()
	defer tun.ListensMutex.Unlock()
	for id, _ := range tun.Listens {
		tun.closeListen(id)
	}
}

func (tun *Tunnel) CloseConns() {
	tun.ConnsMutex.Lock()
	defer tun.ConnsMutex.Unlock()
	for id, _ := range tun.Conns {
		tun.closeConn(id, false)
	}
}

func (tun *Tunnel) ClearListens() {
	tun.ListensMutex.Lock()
	defer tun.ListensMutex.Unlock()

	for id, _ := range tun.Listens {
		tun.closeListen(id)
	}
	tun.Listens = map[uint64]*Listen{}
}

func (tun *Tunnel) ClearConns() {
	tun.ConnsMutex.Lock()
	defer tun.ConnsMutex.Unlock()

	for id, _ := range tun.Conns {
		tun.closeConn(id, false)
	}
	tun.Conns = map[uint64]*Connection{}
}

func (tun *Tunnel) String() string {
	tun.ListensMutex.Lock()
	tun.ConnsMutex.Lock()
	defer tun.ListensMutex.Unlock()
	defer tun.ConnsMutex.Unlock()

	return fmt.Sprintf("{Name:%v, RemoteAddr:%v, Listens:%v, Conns:%v, Error:%v}", tun.Name, tun.RemoteAddr, tun.Listens, tun.Conns, tun.Error)
}
