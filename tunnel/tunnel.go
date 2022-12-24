package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

const (
	PACK_SIZE   = 1024
	BUFFER_SIZE = 4 * 1024
)

type Tunnel struct {
	Direction uint8
	Addr      string

	TunConn   net.Conn
	InBuffer  []byte
	OutBuffer []byte

	ConnsCnt uint64
	Conns    map[uint64]net.Conn

	// listener

	// client

}

func NewTunnel(direction uint8, addr string, tunConn net.Conn) *Tunnel {
	tun := &Tunnel{
		Direction: direction,
		Addr:      addr,

		TunConn:   tunConn,
		InBuffer:  make([]byte, BUFFER_SIZE),
		OutBuffer: make([]byte, BUFFER_SIZE),

		ConnsCnt: 0,
		Conns:    map[uint64]net.Conn{},
	}

	return tun
}

func (tun *Tunnel) ReadMsg() (Msg, error) {
	if _, err := io.ReadFull(tun.TunConn, tun.InBuffer[:5]); err != nil {
		return nil, err
	}

	_, ln := MsgType(tun.InBuffer[0]), binary.LittleEndian.Uint32(tun.InBuffer[1:5])
	if ln+5 > uint32(len(tun.InBuffer)) {
		return nil, fmt.Errorf("msg too big %v", ln)
	}

	if _, err := io.ReadFull(tun.TunConn, tun.InBuffer[5:5+ln]); err != nil {
		return nil, err
	}

	msg, err := deserialize(tun.InBuffer)
	return msg, err
}

func (tun *Tunnel) WriteMsg(msg Msg) error {
	n, err := serialize(msg, tun.OutBuffer)
	if err != nil {
		return err
	}

	_, err = tun.TunConn.Write(tun.OutBuffer[:n])
	return err
}

func (tun *Tunnel) NewId() uint64 {
	tun.ConnsCnt++
	return tun.ConnsCnt
}

func (tun *Tunnel) RunListener() error {
	listen, err := net.Listen("tcp", tun.Addr)
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		tun.OpenConn(conn)
	}
}

// open new connection
func (tun *Tunnel) OpenConn(conn net.Conn) {
	id := tun.NewId()
	tun.Conns[id] = conn
	go tun.ConnectHandler(conn, id)
}

// close connection
func (tun *Tunnel) CloseConn(id uint64) {
	delete(tun.Conns, id)
}

// conn -> tunnel
func (tun *Tunnel) ConnectHandler(conn net.Conn, id uint64) {
	buf := make([]uint8, PACK_SIZE)
	msg := &MsgPack{
		Id: id,
	}

	n, err := 0, error(nil)
	for err == nil {
		if n, err = conn.Read(buf); n > 0 && err == nil {
			msg.Data = buf[:n]
			err = tun.WriteMsg(msg)
		}
	}

	tun.CloseConn(id)
}
