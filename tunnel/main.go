package main

import (
	"flag"
	"log"
	"strings"
)

var D string
var LParas string
var RParas string
var Addr string
var Password string

func main() {
	parseFlag()

	if D == "client" { // connect to server
		tun := &TunnelProxy{}
		copy(tun.Password[:], []uint8(Password))
		tun.Addr = Addr

		if LParas != "" {
			ss := strings.Split(LParas, ":")
			if len(ss) != 4 {
				log.Fatal("error paras %v", LParas)
			}

			tun.LocalAddr = ss[0] + ":" + ss[1]
			tun.RemoteAddr = ss[2] + ":" + ss[3]

		} else if RParas != "" {
			ss := strings.Split(RParas, ":")
			if len(ss) != 4 {
				log.Fatal("error paras %v", RParas)
			}

			tun.RemoteAddr = ss[0] + ":" + ss[1]
			tun.LocalAddr = ss[2] + ":" + ss[3]
		}

		tun.RunClient()

	} else { // run as server
		tun := &TunnelProxy{}
		copy(tun.Password[:], []uint8(Password))
		tun.Addr = Addr
		tun.RunServer()
	}

}

func parseFlag() {
	flag.StringVar(&D, "D", "client", "as daemon")
	flag.StringVar(&LParas, "L", "", "local listener, remote client")
	flag.StringVar(&RParas, "R", "", "local client, remote listener")
	flag.StringVar(&Addr, "H", "127.0.0.1:2333", "ip:port")
	flag.StringVar(&Password, "P", "", "")
}
