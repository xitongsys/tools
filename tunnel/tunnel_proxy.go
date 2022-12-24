package main

import (
	"log"
	"net"
)

type TunnelProxy struct {
	Daemon   bool
	Password [8]uint8

	// listener
	Addr string

	// client
	Direction  string
	LocalAddr  string
	RemoteAddr string
}

func (tp *TunnelProxy) RunServer() {
	listen, err := net.Listen("tcp", tp.Addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listen.Close()

	for {
		tunConn, err := listen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go tp.TunConnHandler(tunConn)
	}
}

func (tp *TunnelProxy) TunConnHandler(tunConn net.Conn) {
	buffer := make([]uint8, BUFFER_SIZE)
	msgi, err := ReadMsg(tunConn, buffer)
	if err != nil {
		return
	}

	if msgi.Type() == TUN {
		msg := msgi.(*MsgTun)
		var addr string
		addr, err = Int2Addr(msg.Addr)
		if err == nil {
			tun := NewTunnel(msg.Role, addr, tunConn)
			go tun.Run()
		}
	}
}

/////////////////////////////////

func (tp *TunnelProxy) RunClient() {
	buffer := make([]uint8, BUFFER_SIZE)
	tunConn, err := net.Dial("tcp", tp.Addr)
	if err != nil {
		log.Fatal(err)
	}
	defer tunConn.Close()

	localRole, remoteRole := 'L', 'C'
	if tp.Direction == "L" {
		localRole, remoteRole = 'L', 'C'

	} else if tp.Direction == "R" {
		localRole, remoteRole = 'C', 'L'

	} else {
		log.Fatal("unknown direction %v", tp.Direction)
	}

	msg := &MsgTun{
		Role:     uint8(remoteRole),
		Password: tp.Password,
	}
	msg.Addr, err = Addr2Int(tp.RemoteAddr)
	if err != nil {
		log.Fatal(err)
	}

	if err = WriteMsg(tunConn, buffer, msg); err != nil {
		log.Fatal(err)
	}

	tun := NewTunnel(uint8(localRole), tp.LocalAddr, tunConn)
	tun.Run()
}
