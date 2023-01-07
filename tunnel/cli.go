package main

import "fmt"

func Cli() {
	var cmd string
	fmt.Scan(&cmd)

	if cmd == "info" {
		fmt.Println(tp.ToString())

	} else if cmd == "listen" {
		var tunName, direction string
		var listenAddr, forwardAddr string
		fmt.Scan(&tunName, &direction, &listenAddr, &forwardAddr)

		if tun := tp.GetTun(tunName); tun != nil {
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

	} else if cmd == "closelisten" {
		var tunName, direction string
		var id uint64
		fmt.Scan(&tunName, &direction, &id)

		if tun := tp.GetTun(tunName); tun != nil {
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

	} else if cmd == "kill" {
		var tunName string
		fmt.Scan(&tunName)
		tp.CloseTun(tunName)
	}
}
