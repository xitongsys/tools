package main

import (
	"fmt"
	"os"
	"flag"
	"net"
	"log"
	"io"
)

var ip, port string
var proxy string
var protocol string


func main() {
	flag.StringVar(&protocol, "X", "connect", "protocol")
	flag.StringVar(&proxy, "x", "127.0.0.1:3128", "http proxy address")
	flag.StringVar(&ip, "ip", "127.0.0.1:3128", "http proxy address")
	flag.StringVar(&port, "port", "127.0.0.1:3128", "http proxy address")
	flag.Parse()


	conn, err := net.Dial("tcp", proxy)	
	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()


	s := fmt.Sprintf("%s %s:%s HTTP/1.1\n\n", protocol,ip,port)
	if _, err := conn.Write([]byte(s)); err != nil {
		log.Println(err)
		return
	}

	buf := make([]byte,1024)
	if _, err := conn.Read(buf); err != nil {
		log.Println(err)
		return
	}


	go func(){
		if _, err := io.Copy(conn, os.Stdin); err != nil {
			log.Println(err)
			return
		}
	}()


	if _, err := io.Copy(os.Stdout, conn); err != nil {
		log.Println(err)
	}

}
