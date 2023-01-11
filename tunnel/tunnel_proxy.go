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

	Tunnels map[string]*Tunnel
	// tun created by self should be keep retrying if error happened
	RetryTunnels map[string]*Tunnel
	TunnelsMutex sync.Mutex
}

func NewTunelProxy(addr string, password string) *TunnelProxy {
	tp := &TunnelProxy{}
	copy(tp.Password[:], []uint8(password))
	tp.CipherBlock = Password2Cipher(password)
	tp.Addr = addr
	tp.Tunnels = map[string]*Tunnel{}
	tp.RetryTunnels = map[string]*Tunnel{}
	tp.Tunnels["local"] = nil

	return tp
}

func (tp *TunnelProxy) Run() {
	// start server
	if tp.Addr != "" {
		go func() {
			listen, err := net.Listen("tcp", tp.Addr)
			if err != nil {
				log.Fatal(err)
			}
			defer listen.Close()

			for {
				conn, err := listen.Accept()
				if err != nil {
					Logger(ERRO, "accept error: %v\n", err)
					continue
				}
				go tp.ConnHandler(conn)
			}
		}()
	}

	// clean job
	go func() {
		for {
			tp.CleanTun()
			time.Sleep(3 * time.Second)
		}
	}()
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
		password := ByteArrayToString(msg.Password[:])
		tun := NewTunnel(name, tunConn.RemoteAddr().String(), password, tunConn)
		Logger(INFO, "new tunnel: %v, %v\n", name, tunConn.RemoteAddr())

		tp.TunnelsMutex.Lock()
		defer tp.TunnelsMutex.Unlock()

		if _, ok := tp.Tunnels[name]; ok {
			tunConn.Close()
			Logger(WARN, "duplicated tun: %v %v\n", name, tun.RemoteAddr)
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
			tun := NewTunnel(name, addr, password, conn)
			go tun.Run()

			tp.TunnelsMutex.Lock()
			tp.Tunnels[name] = tun
			tp.RetryTunnels[name] = tun
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

	delete(tp.RetryTunnels, name)
}

func (tp *TunnelProxy) CleanTun() {
	delete_names := []string{}

	tp.TunnelsMutex.Lock()
	for name, tun := range tp.Tunnels {
		if tun == nil || tun.Error != nil {
			if _, ok := tp.RetryTunnels[name]; !ok {
				delete_names = append(delete_names, name)
				if tun != nil {
					Logger(WARN, "tun closed: %v %v %v", tun.Name, tun.RemoteAddr, tun.Error)
				}
			}
		}
	}

	for _, name := range delete_names {
		tun := tp.Tunnels[name]
		if tun != nil {
			tun.ClearConns()
			tun.ClearListens()
			Logger(WARN, "clean tun %v", name)
		}
		delete(tp.Tunnels, name)
	}
	tp.TunnelsMutex.Unlock()

	for _, tun := range tp.RetryTunnels {
		if tun.Error != nil {
			err := tp.OpenTun(tun.RemoteAddr, tun.Name, tun.Password)
			Logger(WARN, "retry tun %v %v %v", tun.Name, tun.RemoteAddr, err)

			tun.CloseListens()
			tun.ClearConns()

			if err == nil {
				for _, listen := range tun.Listens {
					err = tun.OpenListen(listen.Id, listen.ListenAddr, listen.ForwardAddr)
					Logger(WARN, "retry listen %v %v %v %v", listen.Id, listen.ListenAddr, listen.ForwardAddr, err)
				}
			}
		}
	}
}

func (tp *TunnelProxy) String() string {
	tp.TunnelsMutex.Lock()
	defer tp.TunnelsMutex.Unlock()

	tuns := []string{}
	for name, _ := range tp.Tunnels {
		tuns = append(tuns, name)
	}

	retryTuns := []string{}
	for name, _ := range tp.RetryTunnels {
		retryTuns = append(retryTuns, name)
	}

	res := fmt.Sprintf(`{
		Addr: %v,
		Password: %v,
		Tunnels: %v,	
		RetryTunnels: %v,
	}`, tp.Addr, string(tp.Password[:]), tuns, retryTuns)

	return res
}

/////////////////////////////////
