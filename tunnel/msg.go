package main

import (
	"encoding/binary"
	"fmt"
)

type MsgType uint8

const (
	DEFAULT MsgType = iota
	TUN
	CONN
	CLOSE
	PACK
)

type Msg interface {
	Type() MsgType
	Len() uint32
}

//////////////////

// open new tunnel
type MsgTun struct {
	Role     uint8
	Addr     uint64
	Password [8]uint8
}

func (msgtun *MsgTun) Len() uint32 {
	return 1 + 8 + 8
}

func (msgtun *MsgTun) Type() MsgType {
	return TUN
}

//////////////////////

// new connection
type MsgConn struct {
	Id uint64
}

func (msgconn *MsgConn) Len() uint32 {
	return 8
}

func (msgconn *MsgConn) Type() MsgType {
	return CONN
}

//////////////////////

// close connection
type MsgClose struct {
	Id uint64
}

func (msgclose *MsgClose) Len() uint32 {
	return 8
}

func (msgclose *MsgClose) Type() MsgType {
	return CLOSE
}

/////////////////////

// network package
type MsgPack struct {
	Id   uint64
	Data []byte
}

func (msgpack *MsgPack) Len() uint32 {
	return 8 + uint32(len(msgpack.Data))
}

func (msgpack *MsgPack) Type() MsgType {
	return PACK
}

/////////////////////

func serialize(msgi Msg, buf []uint8) (uint32, error) {
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

		buf[offset] = msg.Role
		offset++

		binary.LittleEndian.PutUint64(buf[offset:], msg.Addr)
		offset += 8

		copy(buf[offset:], msg.Password[:])
		offset += 8

	} else if tp == CONN {
		msg := msgi.(*MsgConn)

		binary.LittleEndian.PutUint64(buf[offset:], msg.Id)
		offset += 8

	} else if tp == CLOSE {
		msg := msgi.(*MsgClose)

		binary.LittleEndian.PutUint64(buf[offset:], msg.Id)
		offset += 8

	} else if tp == PACK {
		msg := msgi.(*MsgPack)

		binary.LittleEndian.PutUint64(buf[offset:], msg.Id)
		offset += 8

		copy(buf[offset:], msg.Data)
		offset += len(msg.Data)
	}

	return uint32(offset), nil

}

func deserialize(buf []uint8) (Msg, error) {
	ln := uint32(len(buf))
	tp, offset := DEFAULT, 0
	msgi := Msg(nil)

	if ln < 5 {
		goto ERROR
	}
	tp = MsgType(buf[0])
	offset += 1

	ln = binary.LittleEndian.Uint32(buf[offset:])
	offset += 4

	if ln+5 > uint32(len(buf)) {
		goto ERROR
	}

	if tp == TUN {
		msg := &MsgTun{}

		if offset+1 > len(buf) {
			goto ERROR
		}
		msg.Role = buf[offset]
		offset += 1

		if offset+8 > len(buf) {
			goto ERROR
		}
		msg.Addr = binary.LittleEndian.Uint64(buf[offset:])
		offset += 8

		if offset+8 > len(buf) {
			goto ERROR
		}
		copy(msg.Password[:], buf[offset:])
		offset += 8

		msgi = msg

	} else if tp == CONN {
		msg := &MsgConn{}

		if offset+8 > len(buf) {
			goto ERROR
		}
		msg.Id = binary.LittleEndian.Uint64(buf[offset:])
		offset += 8

		msgi = msg

	} else if tp == CLOSE {
		msg := &MsgClose{}

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

		if offset+int(ln) > len(buf) || int(ln) < 0 {
			goto ERROR
		}
		msg.Data = make([]byte, ln)
		copy(msg.Data, buf[offset:offset+int(ln)])
		offset += int(ln)

		msgi = msg
	}

	return msgi, nil

ERROR:
	return nil, fmt.Errorf("error msg")
}

//////////////////////
