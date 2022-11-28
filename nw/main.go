package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"flag"
	"io"
	"log"
	"net"
	"os"
)

const BUFSIZE = 1024 * 16

var src string
var dst string
var srcPwd string
var dstPwd string
var protocol string

func main() {
	parseFlag()
	listen, err := net.Listen(protocol, dst)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleNewConn(conn)
	}
}

func parseFlag() {
	flag.StringVar(&src, "src", "", "source address ip:port")
	flag.StringVar(&dst, "dst", "", "destination address ip:port")
	flag.StringVar(&srcPwd, "srcpwd", "", "AES password for soruce")
	flag.StringVar(&dstPwd, "dstpwd", "", "AES password for destination")
	flag.StringVar(&protocol, "pro", "tcp", "the network protocol tcp/udp")
	flag.Parse()

	if srcPwd != "" {
		for len(srcPwd) < 16 {
			srcPwd = srcPwd + "a"
		}
		srcPwd = srcPwd[:16]
	}

	if dstPwd != "" {
		for len(dstPwd) < 16 {
			dstPwd = dstPwd + "a"
		}
		dstPwd = dstPwd[:16]
	}
}

func handleNewConn(dstConn net.Conn) {
	srcConn, err := net.Dial(protocol, src)
	if err != nil {
		log.Println(err)
		return
	}

	// dst -> src
	go copy(srcConn, srcPwd, dstConn, dstPwd)
	// src -> dst
	go copy(dstConn, dstPwd, srcConn, srcPwd)
}

func encrypt(dst []byte, src []byte, block cipher.Block) {
	n, nb := len(dst), block.BlockSize()
	for i := 0; i*nb < n; i++ {
		bgn, end := i*nb, i*nb+nb
		block.Encrypt(dst[bgn:end], src[bgn:end])
	}
}

func decrypt(dst []byte, src []byte, block cipher.Block) {
	n, nb := len(dst), block.BlockSize()
	for i := 0; i*nb < n; i++ {
		bgn, end := i*nb, i*nb+nb
		block.Decrypt(dst[bgn:end], src[bgn:end])
	}
}

func copy(w io.Writer, wPwd string, r io.Reader, rPwd string) {
	if rPwd != "" && wPwd != "" {
		rBuf0, wBuf0, mBuf0 := make([]byte, BUFSIZE), make([]byte, BUFSIZE), make([]byte, BUFSIZE)
		sBuf := make([]byte, 4)
		rc, err := aes.NewCipher([]byte(rPwd))
		if err != nil {
			log.Println(err)
			return
		}

		wc, err := aes.NewCipher([]byte(wPwd))
		if err != nil {
			log.Println(err)
			return
		}

		for {
			if _, err := io.ReadFull(r, sBuf); err != nil {
				log.Println(err)
				return
			}

			len := binary.LittleEndian.Uint32(sBuf)
			// paddling length
			plen := ((len + 15) >> 4) << 4

			if plen > BUFSIZE {
				log.Println(err)
				return
			}

			rBuf, wBuf, mBuf := rBuf0[:plen], wBuf0[:plen], mBuf0[:plen]
			if _, err := io.ReadFull(r, rBuf); err != nil {
				log.Println(err)
				return
			}

			decrypt(mBuf, rBuf, rc)
			encrypt(wBuf, mBuf, wc)

			if _, err := w.Write(sBuf); err != nil {
				log.Println(err)
				return
			}

			if _, err := w.Write(wBuf); err != nil {
				log.Println(err)
				return
			}
		}

	} else if rPwd != "" && wPwd == "" {
		rBuf0, wBuf0 := make([]byte, BUFSIZE), make([]byte, BUFSIZE)
		sBuf := make([]byte, 4)
		rc, err := aes.NewCipher([]byte(rPwd))
		if err != nil {
			log.Println(err)
			return
		}

		for {
			if _, err := io.ReadFull(r, sBuf); err != nil {
				log.Println(err)
				return
			}

			len := binary.LittleEndian.Uint32(sBuf)
			// paddling length
			plen := ((len + 15) >> 4) << 4
			if plen > BUFSIZE {
				log.Println(err)
				return
			}

			rBuf, wBuf := rBuf0[:plen], wBuf0[:plen]
			if _, err := io.ReadFull(r, rBuf); err != nil {
				log.Println(err)
				return
			}

			decrypt(wBuf, rBuf, rc)

			if _, err := w.Write(wBuf[:len]); err != nil {
				log.Println(err)
				return
			}
		}

	} else if rPwd == "" && wPwd != "" {
		rBuf0, wBuf0 := make([]byte, BUFSIZE), make([]byte, BUFSIZE)
		sBuf := make([]byte, 4)
		wc, err := aes.NewCipher([]byte(wPwd))
		if err != nil {
			log.Println(err)
			return
		}

		for {
			len, err := r.Read(rBuf0)
			if err != nil {
				log.Println(err)
				return
			}

			// paddling length
			plen := ((len + 15) >> 4) << 4
			if uint32(plen) > BUFSIZE {
				log.Println(err)
				return
			}

			rBuf, wBuf := rBuf0[:plen], wBuf0[:plen]

			binary.LittleEndian.PutUint32(sBuf, uint32(len))

			encrypt(wBuf, rBuf, wc)

			if _, err := w.Write(sBuf); err != nil {
				log.Println(err)
				return
			}

			if _, err := w.Write(wBuf); err != nil {
				log.Println(err)
				return
			}
		}

	} else {
		if _, err := io.Copy(w, r); err != nil {
			log.Println(err)
			return
		}
	}
}
