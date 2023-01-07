package main

import (
	"crypto/aes"
	"crypto/cipher"
	"log"
	"net"
	"sync"
)

type TunnelProxy struct {
	Password    [16]uint8
	Addr        string
	CipherBlock cipher.Block

	Tunnels      map[string]*Tunnel
	TunnelsMutex sync.Mutex
}

func NewTunelProxy(addr string, password string) *TunnelProxy {
	tp := &TunnelProxy{}
	for i := 0; i < 16 && i < len(password); i++ {
		tp.Password[i] = password[i]
	}

	if len(password) > 0 {
		tp.CipherBlock, _ = aes.NewCipher(tp.Password[:])
	}

	tp.Addr = addr
	tp.Tunnels = map[string]*Tunnel{}

	return tp
}

func (tp *TunnelProxy) Run() {
	listen, err := net.Listen("tcp", tp.Addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go tp.ConnHandler(conn)
	}
}

func (tp *TunnelProxy) ConnHandler(tunConn net.Conn) {
	buffer := make([]uint8, BUFFER_SIZE)
	msgi, err := ReadMsg(tunConn, buffer, tp.CipherBlock)
	if err != nil {
		tunConn.Close()
		log.Println(err)
		return
	}

	if msgi.Type() == TUN {
		msg := msgi.(*MsgTun)

		if msg.Password != tp.Password {
			tunConn.Close()
			return
		}

		name := string(msg.Name[:])

		tun := NewTunnel(name, tunConn, tp.CipherBlock)
		log.Printf("new tunnel: %v, %v", name, tunConn.RemoteAddr())

		tp.TunnelsMutex.Lock()
		defer tp.TunnelsMutex.Unlock()

		if _, ok := tp.Tunnels[name]; ok {
			tunConn.Close()
			return
		}

		tp.Tunnels[name] = tun
		go tun.Run()
	}
}

/////////////////////////////////
