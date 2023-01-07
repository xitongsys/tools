package main

import (
	"io"
	"log"
)

type LogLevel uint8

const (
	DEBG LogLevel = iota
	INFO
	WARN
	ERRO
)

var Prefixs []string = []string{"DEBG", "INFO", "WARN", "ERRO"}

var CurrentLogLevel LogLevel = DEBG
var logger *log.Logger

func InitLog(file io.Writer) {
	logger = log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func Logger(logLevel LogLevel, fmt string, v ...any) {
	if logLevel >= CurrentLogLevel {
		prefix := Prefixs[logLevel]
		logger.Printf(prefix+fmt, v...)
	}
}
