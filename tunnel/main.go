package main

import (
	"flag"
	"log"
	"os"
)

var Addr string
var Password string
var LogFile string

var TP *TunnelProxy

func main() {
	parseFlag()

	// log
	if LogFile == "" {
		InitLog(os.Stderr)

	} else {
		file, err := os.OpenFile(LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		InitLog(file)
	}

	TP = NewTunelProxy(Addr, Password)
	go TP.Run()

	Cli()
}

func parseFlag() {
	flag.StringVar(&LogFile, "log", "", "")
	flag.StringVar(&Addr, "addr", "0.0.0.0:23", "")
	flag.StringVar(&Password, "pwd", "12345", "")

	flag.Parse()
}
