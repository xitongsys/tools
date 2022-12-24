package main

import (
	"flag"
	"log"
	"strings"
)

var Role string
var LParas string
var RParas string
var Addr string
var Password string

func main() {
	parseFlag()

	if Role == "client" { // connect to server
		tun := NewTunelProxy(Password, Addr)

		if LParas != "" {
			ss := strings.Split(LParas, ":")
			if len(ss) != 4 {
				log.Fatalf("error paras %v", LParas)
			}

			tun.LocalAddr = ss[0] + ":" + ss[1]
			tun.RemoteAddr = ss[2] + ":" + ss[3]
			tun.Direction = "L"

		} else if RParas != "" {
			ss := strings.Split(RParas, ":")
			if len(ss) != 4 {
				log.Fatalf("error paras %v", RParas)
			}

			tun.RemoteAddr = ss[0] + ":" + ss[1]
			tun.LocalAddr = ss[2] + ":" + ss[3]
			tun.Direction = "R"
		}

		tun.RunClient()

	} else { // run as server
		tun := NewTunelProxy(Password, Addr)
		tun.RunServer()
	}

}

func parseFlag() {
	flag.StringVar(&Role, "D", "client", "as daemon")
	flag.StringVar(&LParas, "L", "", "local listener, remote client")
	flag.StringVar(&RParas, "R", "", "local client, remote listener")
	flag.StringVar(&Addr, "H", "127.0.0.1:2333", "ip:port")
	flag.StringVar(&Password, "P", "", "")
	flag.Parse()
}
