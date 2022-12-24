package main

import (
	"log"
	"net"
)

type TunnelProxy struct {
	Role string

	Addr     string
	Password [8]uint8

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
			tun := NewTunnel(msg.Direction, addr, tunConn)
			go tun.Run()
		}
	}
}

/////////////////////////////////

func (tp *TunnelProxy) RunClient() {

}
