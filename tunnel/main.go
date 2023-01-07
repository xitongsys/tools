package main

import (
	"flag"
)

var TunsParas string

func main() {
	parseFlag()

}

func parseFlag() {
	flag.StringVar(&TunsParas, "T", "127.0.0.1:22:12345,127.0.0.1:22:12345", "")
	flag.Parse()
}
