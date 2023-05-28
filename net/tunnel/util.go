package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"net"
	"strconv"
	"strings"
)

// "127.0.0.1:22" -> uint64
func Addr2Int(addr string) (uint64, error) {
	if ss := strings.Split(addr, ":"); len(ss) == 2 {
		if port, err := strconv.Atoi(ss[1]); err == nil {
			if ips := strings.Split(ss[0], "."); len(ips) == 4 {
				if a, err := strconv.Atoi(ips[0]); err == nil {
					if b, err := strconv.Atoi(ips[1]); err == nil {
						if c, err := strconv.Atoi(ips[2]); err == nil {
							if d, err := strconv.Atoi(ips[3]); err == nil {
								res := uint64(0)
								res |= uint64(a) << (32 + 24)
								res |= uint64(b) << (32 + 16)
								res |= uint64(c) << (32 + 8)
								res |= uint64(d) << (32 + 0)
								res |= uint64(port)

								return res, nil
							}
						}
					}
				}
			}
		}
	}
	return 0, fmt.Errorf("illegal addr: %v", addr)
}

// uint64 -> "127.0.0.1:22"
func Int2Addr(ai uint64) (string, error) {
	a, b, c, d := (ai>>(32+24))&0xff, (ai>>(32+16))&0xff, (ai>>(32+8))&0xff, (ai>>(32+0))&0xff
	port := ai & 0xffffffff

	if port > 0 && port < (1<<16) {
		return fmt.Sprintf("%v.%v.%v.%v:%v", a, b, c, d, port), nil
	}

	return "", fmt.Errorf("illegal addr: %v", ai)
}

// read msg from conn
func ReadMsg(conn net.Conn, buffer []uint8, cipherBlock cipher.Block) (Msg, error) {
	if _, err := io.ReadFull(conn, buffer[:5]); err != nil {
		return nil, err
	}

	_, ln := MsgType(buffer[0]), binary.LittleEndian.Uint32(buffer[1:5])
	if ln+5 > uint32(len(buffer)) {
		return nil, fmt.Errorf("msg too big %v", ln)
	}

	if _, err := io.ReadFull(conn, buffer[5:5+ln]); err != nil {
		return nil, err
	}

	msg, err := deserialize(buffer, cipherBlock)

	return msg, err
}

// write msg to conn
func WriteMsg(conn net.Conn, buffer []uint8, msg Msg, cipherBlock cipher.Block) error {
	n, err := serialize(msg, buffer, cipherBlock)
	if err != nil {
		return err
	}
	_, err = conn.Write(buffer[:n])
	return err
}

func PaddlingLen(ln uint32) uint32 {
	return ((ln + 15) >> 4) << 4
}

func Paddling(buffer []uint8) {
	for i := 0; i < len(buffer); i++ {
		buffer[i] = 0
	}
}

func Encrypt(dst []byte, src []byte, block cipher.Block) int {
	nsrc, nb := len(src), block.BlockSize()
	for i := 0; i*nb < nsrc; i++ {
		bgn, end := i*nb, i*nb+nb
		block.Encrypt(dst[bgn:end], src[bgn:end])
	}
	return nsrc
}

func Decrypt(dst []byte, src []byte, block cipher.Block) int {
	nsrc, nb := len(src), block.BlockSize()
	for i := 0; i*nb < nsrc; i++ {
		bgn, end := i*nb, i*nb+nb
		block.Decrypt(dst[bgn:end], src[bgn:end])
	}

	return nsrc
}

func UUID() uint64 {
	id := uint64(rand.Int63())
	return id
}

func Password2Cipher(password string) cipher.Block {
	bs := [16]uint8{}
	for i := 0; i < len(bs) && i < len(password); i++ {
		bs[i] = password[i]
	}

	cipherBlock, _ := aes.NewCipher([]uint8(password))
	return cipherBlock
}

// '\0' as string end
func StringByteLen(bs []uint8) (l int) {
	for l = 0; l < len(bs) && bs[l] != 0; l++ {
	}
	return
}

func ByteArrayToString(bs []uint8) string {
	return string(bs[:StringByteLen(bs)])
}
