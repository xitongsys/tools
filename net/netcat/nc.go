package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var hostname string
var port int

//-x 127.0.0.1:3128
var proxy_address string

//-X connect
var proxy_protocol string

//-i interval
var interval int

//-s 127.0.0.1
var source_ip_addr string

//-p 33
var source_port int

//-l
var is_listen bool

func parse_paras() (err error) {
	non_flag_paras := []string{}
	bool_flag_paras := "46DdhklnrStUuvzC"

	n := len(os.Args)
	i := 1
	for i < n {
		p := os.Args[i]
		if p[0] == '-' {
			p = p[1:]

			// bool flag params
			if strings.Contains(bool_flag_paras, p) {
				if p == "l" {
					is_listen = true
					if source_ip_addr != "" || source_port != 0 {
						return fmt.Errorf("-l conflicts with source_ip_addr/source_port")
					}

				} else {
					return fmt.Errorf("unknown params %s", p)
				}

			} else { // flag params

				if i+1 >= n {
					return fmt.Errorf("lack params")
				}

				np := os.Args[i+1]

				if p == "x" {
					proxy_address = np

				} else if p == "X" {
					proxy_protocol = np
					if np != "4" && np != "5" && np != "connect" {
						return fmt.Errorf("unknown proxy_protocol %s", np)
					}

				} else if p == "s" {
					if is_listen {
						return fmt.Errorf("-l conflicts with source_ip_addr")
					}
					source_ip_addr = np

				} else if p == "p" {
					if is_listen {
						return fmt.Errorf("-l conflicts with source_port")
					}
					if source_port, err = strconv.Atoi(np); err != nil {
						return err
					}

				} else if p == "i" {
					if interval, err = strconv.Atoi(np); err != nil {
						return err
					}

				} else {
					return fmt.Errorf("unknown params %s", p)
				}

				i += 1
			}

		} else {
			non_flag_paras = append(non_flag_paras, p)
		}

		i += 1
	}

	if len(non_flag_paras) != 2 {
		return fmt.Errorf("lack hostname and port")
	}

	hostname = non_flag_paras[0]

	if port, err = strconv.Atoi(non_flag_paras[1]); err != nil {
		return err
	}

	return
}

func listen_mode() {
	log.Println("listen mode not support")
}

func connect_mode() {
	if conn, err := open_conn(); err != nil {
		log.Println(err)
		return
	} else {
		go func() {
			if _, err := io.Copy(conn, os.Stdin); err != nil {
				log.Println(err)
				return
			}
		}()

		if _, err := io.Copy(os.Stdout, conn); err != nil {
			log.Println(err)
		}
	}
}

func open_conn() (conn io.ReadWriter, err error) {
	if proxy_protocol == "" { // direct connect
		addr := fmt.Sprintf("%s:%d", hostname, port)
		return net.Dial("tcp", addr)

	} else if proxy_protocol == "connect" { // http proxy
		if conn, err := net.Dial("tcp", proxy_address); err != nil {
			return nil, err
		} else {
			req := fmt.Sprintf("CONNECT %s:%d HTTP/1.1\n\n", hostname, port)
			if _, err := conn.Write([]byte(req)); err != nil {
				log.Println(err)
				return nil, err
			}

			buf := make([]byte, 1024)
			is_first_line := true
			for {
				i := 0
				for i = 0; i < len(buf); i++ {
					if _, err = conn.Read(buf[i : i+1]); err != nil {
						return nil, fmt.Errorf("proxy error response")
					}

					if buf[i] == '\n' {
						break
					}
				}

				if i >= len(buf) {
					return nil, fmt.Errorf("proxy error response: msg too long")
				}

				if is_first_line {
					line := string(buf[:i])
					tokens := strings.Split(line, " ")
					if len(tokens) < 3 || tokens[1] != "200" {
						return nil, fmt.Errorf("proxy error response: %s", line)
					}
					is_first_line = false
				}

				if i == 1 {
					break
				}
			}
			return conn, err
		}

	} else if proxy_protocol == "4" { // SOCKS4
		return nil, nil
	} else if proxy_protocol == "5" { // SOCKS5
		return nil, nil
	} else {
		return nil, fmt.Errorf("unknown protocol %s", proxy_protocol)
	}
}

func main() {
	if err := parse_paras(); err != nil {
		log.Println(err)
		return
	}

	if is_listen {
		listen_mode()
	} else {
		connect_mode()
	}
}
