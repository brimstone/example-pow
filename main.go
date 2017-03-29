package main

import (
	"bytes"
	"crypto/sha256"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"
	"syscall"
)

/*
func check(msg string, bitlen int) bool {
	return bytes.HasPrefix(hashed[:], prefix) && hashed[offset]&partial == 0
}
*/

func run(msg string, start int64, inc int64, bitlen int, done chan bool) int64 {
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
		pows = strconv.FormatInt(pow, 10)
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

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	msg := flag.Arg(0)

	done := make(chan bool)
	c := int64(runtime.NumCPU())
	var i int64

	syscall.Setpriority(syscall.PRIO_PROCESS, 0, 19)

	log.Println("Starting", c, "processes")

	for i = 0; i <= c; i++ {
		go run(msg, i, c, 24, done)
	}

	<-done
}
