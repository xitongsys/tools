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
	TP.Run()

	Cli()
}

func parseFlag() {
	flag.StringVar(&LogFile, "L", "", "log file, default is stdout/stderr")
	flag.StringVar(&Addr, "D", "", "run as daemon with address")
	flag.StringVar(&Password, "P", "12345", "password")

	flag.Parse()
}
