package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"syscall"
)

/*
func check(msg string, bitlen int) bool {
	return bytes.HasPrefix(hashed[:], prefix) && hashed[offset]&partial == 0
}
*/

func run(msg string, start int, inc int, bitlen int, done chan bool) int {
	pow := start
	bits := strconv.Itoa(bitlen)
	header := msg + " <pow:" + bits + ":"
	offset := bitlen / 8
	prefix := bytes.Repeat([]byte{0xff}, offset)
	shift := uint(bitlen % 8)
	partiali := (255 << (8 - shift)) % 256
	partial := byte(partiali)
	var pows string
	var hashed [32]byte
	for {
		pows = fmt.Sprintf("%d", pow)
		hashed = sha256.Sum256([]byte(header + pows + ">"))
		if bytes.HasPrefix(hashed[:], prefix) && hashed[offset]&partial == partial {
			break
		}

		pow = pow + inc
	}
	fmt.Printf("%s %d %x\n", header+pows+">", pow, hashed[:])
	done <- true
	return pow
}

func main() {
	msg := os.Args[1]

	done := make(chan bool)
	c := runtime.NumCPU()

	syscall.Setpriority(syscall.PRIO_PROCESS, 0, 19)

	log.Println("Starting", c, "processes")

	for i := 0; i <= c; i++ {
		go run(msg, i, c, 32, done)
	}

	<-done
	os.Exit(0)

}
