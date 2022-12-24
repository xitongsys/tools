package main

type MsgType int8

const (
	TUN MsgType = 0
	CONN
	CLOSE
	PACK
)

type MsgHeader struct {
	Len  int32
	Type MsgType
}

type MsgTun struct {
	MsgHeader
	Direction  byte
	LocalAddr  string
	RemoteAddr string
	Password   string
}

type MsgConn struct {
	MsgHeader
	Id int64
}

type MsgClose struct {
	MsgHeader
	Id int64
}

type MsgPack struct {
	MsgHeader
	Id   int64
	Data []byte
}
