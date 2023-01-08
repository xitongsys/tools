package main

import (
	"crypto/cipher"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
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
	copy(tp.Password[:], []uint8(password))
	tp.CipherBlock = Password2Cipher(password)
	tp.Addr = addr
	tp.Tunnels = map[string]*Tunnel{}
	tp.Tunnels["local"] = nil

	return tp
}

func (tp *TunnelProxy) Run() {
	listen, err := net.Listen("tcp", tp.Addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listen.Close()

	// clean job
	go func() {
		for {
			tp.CleanTun()
			time.Sleep(10 * time.Second)
		}
	}()

	///////
	for {
		conn, err := listen.Accept()
		if err != nil {
			Logger(ERRO, "accept error: %v\n", err)
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
		Logger(ERRO, "%v", err)
		return
	}

	if msgi.Type() == TUN {
		msg := msgi.(*MsgTun)

		if msg.Password != tp.Password {
			tunConn.Close()
			return
		}

		name := ByteArrayToString(msg.Name[:])
		tun := NewTunnel(name, tunConn, tp.CipherBlock)
		Logger(INFO, "new tunnel: %v, %v\n", name, tunConn.RemoteAddr())

		tp.TunnelsMutex.Lock()
		defer tp.TunnelsMutex.Unlock()

		if _, ok := tp.Tunnels[name]; ok {
			tunConn.Close()
			Logger(WARN, "duplicated tun name: %v\n", name)
			return
		}

		tp.Tunnels[name] = tun
		go tun.Run()
	}
}

func (tp *TunnelProxy) OpenTun(addr string, name string, password string) error {
	buffer := make([]uint8, BUFFER_SIZE)
	if conn, err := net.Dial("tcp", addr); err == nil {
		msg := &MsgTun{}
		copy(msg.Password[:], password)
		copy(msg.Name[:], name)
		cipherBlock := Password2Cipher(password)
		if err = WriteMsg(conn, buffer, msg, cipherBlock); err == nil {
			tun := NewTunnel(name, conn, cipherBlock)
			go tun.Run()

			tp.TunnelsMutex.Lock()
			tp.Tunnels[name] = tun
			tp.TunnelsMutex.Unlock()

		} else {
			return err
		}

		return nil

	} else {
		return err
	}
}

func (tp *TunnelProxy) GetTun(name string) *Tunnel {
	tp.TunnelsMutex.Lock()
	defer tp.TunnelsMutex.Unlock()

	if tun, ok := tp.Tunnels[name]; ok {
		return tun
	}
	return nil
}

func (tp *TunnelProxy) CloseTun(name string) {
	tp.TunnelsMutex.Lock()
	defer tp.TunnelsMutex.Unlock()

	if tun, ok := tp.Tunnels[name]; ok {
		tun.Exit(fmt.Errorf("close"))
		delete(tp.Tunnels, name)
	}
}

func (tp *TunnelProxy) CleanTun() {
	tp.TunnelsMutex.Lock()
	defer tp.TunnelsMutex.Unlock()

	names := []string{}

	for name, tun := range tp.Tunnels {
		if tun == nil || tun.Error != nil {
			names = append(names, name)
		}
	}

	for _, name := range names {
		delete(tp.Tunnels, name)
	}
}

func (tp *TunnelProxy) String() string {
	tp.TunnelsMutex.Lock()
	defer tp.TunnelsMutex.Unlock()

	tuns := []string{}
	for name, _ := range tp.Tunnels {
		tuns = append(tuns, name)
	}

	res := fmt.Sprintf(`{
		Addr: %v,
		Password: %v,
		Tuns: %v,	
	}`, tp.Addr, string(tp.Password[:]), tuns)

	return res
}

/////////////////////////////////
