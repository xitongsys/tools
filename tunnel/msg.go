package main

import (
	"crypto/cipher"
	"encoding/binary"
	"fmt"
)

type MsgType uint8

const (
	DEFAULT MsgType = iota
	TUN

	LISTEN
	CLOSELISTEN

	CONN
	CLOSECONN

	PACK
	CMD
)

type Msg interface {
	Type() MsgType
	Len() uint32
}

//////////////////

// open new tunnel
type MsgTun struct {
	Name     [16]uint8
	Password [16]uint8
}

func (msgtun *MsgTun) Len() uint32 {
	return 16 + 16
}

func (msgtun *MsgTun) Type() MsgType {
	return TUN
}

//////////////////////

// new listen
type MsgListen struct {
	Id          uint64
	ListenAddr  uint64
	ForwardAddr uint64
}

func (msglisten *MsgListen) Len() uint32 {
	return 8 + 8 + 8
}

func (msglisten *MsgListen) Type() MsgType {
	return LISTEN
}

/////////////////

// close connection
type MsgCloseListen struct {
	Id uint64
}

func (msgcloselisten *MsgCloseListen) Len() uint32 {
	return 8
}

func (msgcloselisten *MsgCloseListen) Type() MsgType {
	return CLOSELISTEN
}

/////////////

// new connection
type MsgConn struct {
	Id   uint64
	Addr uint64
}

func (msgconn *MsgConn) Len() uint32 {
	return 8 + 8
}

func (msgconn *MsgConn) Type() MsgType {
	return CONN
}

//////////////////////

// close connection
type MsgCloseConn struct {
	Id uint64
}

func (msgcloseconn *MsgCloseConn) Len() uint32 {
	return 8
}

func (msgcloseconn *MsgCloseConn) Type() MsgType {
	return CLOSECONN
}

/////////////////////

// network package
type MsgPack struct {
	Id      uint64
	DataLen uint32
	Data    []byte
}

func (msgpack *MsgPack) Len() uint32 {
	return 8 + 4 + uint32(len(msgpack.Data))
}

func (msgpack *MsgPack) Type() MsgType {
	return PACK
}

/////////////////////

func serialize(msgi Msg, buf []uint8, cipherBlock cipher.Block) (uint32, error) {
	offset := 0
	tp, ln := msgi.Type(), msgi.Len()

	if uint32(len(buf)) < ln+5 {
		return 0, fmt.Errorf("buf is too small")
	}

	buf[0] = uint8(tp)
	offset += 1

	binary.LittleEndian.PutUint32(buf[offset:], msgi.Len())
	offset += 4

	if tp == TUN {
		msg := msgi.(*MsgTun)

		copy(buf[offset:], msg.Name[:])
		offset += len(msg.Name)

		copy(buf[offset:], msg.Password[:])
		offset += len(msg.Password)

	} else if tp == LISTEN {
		msg := msgi.(*MsgListen)

		binary.LittleEndian.PutUint64(buf[offset:], msg.Id)
		offset += 8

		binary.LittleEndian.PutUint64(buf[offset:], msg.ListenAddr)
		offset += 8

		binary.LittleEndian.PutUint64(buf[offset:], msg.ForwardAddr)
		offset += 8

	} else if tp == CLOSELISTEN {
		msg := msgi.(*MsgCloseListen)

		binary.LittleEndian.PutUint64(buf[offset:], msg.Id)
		offset += 8

	} else if tp == CONN {
		msg := msgi.(*MsgConn)

		binary.LittleEndian.PutUint64(buf[offset:], msg.Id)
		offset += 8

		binary.LittleEndian.PutUint64(buf[offset:], msg.Addr)
		offset += 8

	} else if tp == CLOSECONN {
		msg := msgi.(*MsgCloseConn)

		binary.LittleEndian.PutUint64(buf[offset:], msg.Id)
		offset += 8

	} else if tp == PACK {
		msg := msgi.(*MsgPack)

		binary.LittleEndian.PutUint64(buf[offset:], msg.Id)
		offset += 8

		binary.LittleEndian.PutUint32(buf[offset:], msg.DataLen)
		offset += 4

		if cipherBlock != nil {
			offset += Encrypt(buf[offset:], msg.Data, cipherBlock)

		} else {
			offset += copy(buf[offset:], msg.Data)
		}
	}

	return uint32(offset), nil

}

func deserialize(buf []uint8, cipherBlock cipher.Block) (Msg, error) {
	ln := uint32(len(buf))
	tp, offset := DEFAULT, 0
	msgi := Msg(nil)

	if ln < 5 {
		goto ERROR
	}
	tp = MsgType(buf[0])
	offset += 1

	ln = binary.LittleEndian.Uint32(buf[offset : offset+4])
	offset += 4

	if ln+5 > uint32(len(buf)) {
		goto ERROR
	}

	if tp == TUN {
		msg := &MsgTun{}

		if offset+16 > len(buf) {
			goto ERROR
		}
		copy(msg.Name[:], buf[offset:])
		offset += len(msg.Name)

		if offset+16 > len(buf) {
			goto ERROR
		}
		copy(msg.Password[:], buf[offset:])
		offset += len(msg.Password)

		msgi = msg

	} else if tp == LISTEN {
		msg := &MsgListen{}

		if offset+8 > len(buf) {
			goto ERROR
		}
		msg.Id = binary.LittleEndian.Uint64(buf[offset:])
		offset += 8

		if offset+8 > len(buf) {
			goto ERROR
		}
		msg.ListenAddr = binary.LittleEndian.Uint64(buf[offset:])
		offset += 8

		if offset+8 > len(buf) {
			goto ERROR
		}
		msg.ForwardAddr = binary.LittleEndian.Uint64(buf[offset:])
		offset += 8

		msgi = msg

	} else if tp == CLOSELISTEN {
		msg := &MsgCloseListen{}

		if offset+8 > len(buf) {
			goto ERROR
		}
		msg.Id = binary.LittleEndian.Uint64(buf[offset:])
		offset += 8

		msgi = msg

	} else if tp == CONN {
		msg := &MsgConn{}

		if offset+8 > len(buf) {
			goto ERROR
		}
		msg.Id = binary.LittleEndian.Uint64(buf[offset:])
		offset += 8

		if offset+8 > len(buf) {
			goto ERROR
		}
		msg.Addr = binary.LittleEndian.Uint64(buf[offset:])
		offset += 8

		msgi = msg

	} else if tp == CLOSECONN {
		msg := &MsgCloseConn{}

		if offset+8 > len(buf) {
			goto ERROR
		}
		msg.Id = binary.LittleEndian.Uint64(buf[offset:])
		offset += 8

		msgi = msg

	} else if tp == PACK {
		msg := &MsgPack{}

		if offset+8 > len(buf) {
			goto ERROR
		}
		msg.Id = binary.LittleEndian.Uint64(buf[offset:])
		offset += 8

		if offset+4 > len(buf) {
			goto ERROR
		}
		msg.DataLen = binary.LittleEndian.Uint32(buf[offset:])
		offset += 4

		if offset+int(ln) > len(buf) || int(ln) < 0 {
			goto ERROR
		}
		dataBufferLen := ln + 5 - uint32(offset)
		msg.Data = make([]byte, dataBufferLen)

		if cipherBlock != nil {
			offset += Decrypt(msg.Data, buf[offset:offset+int(dataBufferLen)], cipherBlock)

		} else {
			offset += copy(msg.Data, buf[offset:offset+int(dataBufferLen)])
		}

		msgi = msg
	}

	return msgi, nil

ERROR:
	return nil, fmt.Errorf("error msg")
}

//////////////////////
