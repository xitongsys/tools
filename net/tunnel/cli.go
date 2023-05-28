package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const helpString = `
ls
	show info

ls tunname
	show tun info

open tun tun_name remote_addr remote_password
	open a tun connection named tun_name to remote_addr with remote_password

open listen tun_name direction listen_addr forward_addr
	direction: l(listen on local) or r(listen on remote)
	listen on listen_addr and forward the connection package to forward_addr by the tun
	
close tun tun_name
	close tun

close listen tun_name direction listen_id
	close listen

exit
	exit
`

func checkParasNumber(fields []string, num int) bool {
	if len(fields) < num {
		fmt.Println("illegal paras")
		return false
	}
	return true
}

func Cli() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if !scanner.Scan() {
			time.Sleep(3 * time.Second)
			continue
		}

		cmdline := scanner.Text()
		paras := strings.Fields(cmdline)
		if len(paras) == 0 {
			continue
		}

		cmd := paras[0]
		paras = paras[1:]

		if cmd == "help" {
			fmt.Print(helpString)

		} else if cmd == "ls" {
			if len(paras) == 0 {
				fmt.Printf("%v\n", TP)
			} else {
				name := paras[0]
				tun := TP.GetTun(name)
				if tun != nil {
					fmt.Printf("%v\n", tun)
				} else {
					fmt.Printf("%v not found\n", name)
				}
			}

		} else if cmd == "open" && checkParasNumber(paras, 1) {
			if paras[0] == "tun" && checkParasNumber(paras, 1+3) {
				name, addr, password := paras[1], paras[2], paras[3]
				if err := TP.OpenTun(addr, name, password); err != nil {
					fmt.Println(err)
				}

			} else if paras[0] == "listen" && checkParasNumber(paras, 1+4) {
				tunName, direction, listenAddr, forwardAddr := paras[1], paras[2], paras[3], paras[4]
				if tun := TP.GetTun(tunName); tun != nil {
					if direction == "l" {
						tun.OpenListen(tun.NewId(), listenAddr, forwardAddr)

					} else if direction == "r" {
						listenAddrInt, err1 := Addr2Int(listenAddr)
						forwardAddr, err2 := Addr2Int(forwardAddr)

						if err1 == nil && err2 == nil {
							msg := &MsgListen{
								Id:          tun.NewId(),
								ListenAddr:  listenAddrInt,
								ForwardAddr: forwardAddr,
							}

							tun.WriteMsg(msg)

						} else {
							fmt.Println(err1, err2)
						}

					} else {
						fmt.Printf("unknown direction: %v\n", direction)
					}
				}
			}

		} else if cmd == "close" && checkParasNumber(paras, 1) {
			if paras[0] == "listen" && checkParasNumber(paras, 1+3) {
				tunName, direction, idStr := paras[1], paras[2], paras[3]
				id, _ := strconv.ParseUint(idStr, 10, 64)

				if tun := TP.GetTun(tunName); tun != nil {
					if direction == "l" {
						tun.CloseListen(id)

					} else if direction == "r" {
						msg := &MsgCloseListen{
							Id: id,
						}

						tun.WriteMsg(msg)

					} else {
						fmt.Printf("unknown direction: %v\n", direction)
					}
				}

			} else if paras[0] == "tun" && checkParasNumber(paras, 1+1) {
				tunName := paras[1]
				TP.CloseTun(tunName)
			}

		} else if cmd == "exit" {
			os.Exit(0)

		} else {
			fmt.Printf("unknown cmd: %v\n", cmd)
		}
	}
}
