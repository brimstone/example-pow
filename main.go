package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strconv"
)

/*
func check(msg string, bitlen int) bool {
	return bytes.HasPrefix(hashed[:], prefix) && hashed[offset]&partial == 0
}
*/

func run(msg string, start int, inc int, bitlen int) int {
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
		if bytes.HasPrefix(hashed[:], prefix) && hashed[offset]&partial == 0 {
			break
		}

		pow = pow + inc
	}
	fmt.Printf("%s %d %x\n", header+pows+">", pow, hashed[:])
	return pow
}

func main() {
	msg := "@joeerl n/m figured it out :)"

	run(msg, 0, 1, 24)
}
