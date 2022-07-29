// Xlog v01

package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

var sbox = []int{
	// removed from preview
}

var blowfishSub = []int{
	0x9C30D539,
	0x2AF26013,
	0xC5D1B023,
	0x286085F0,
}

var blowfishP = []int{
	// removed from preview
}

func subDword(x int) int {
	a := (x >> 24) & 0xFF
	b := (x >> 16) & 0xFF
	c := (x >> 8) & 0xFF
	d := x & 0xFF
	y := sbox[d]
	y |= sbox[c] << 8
	y |= sbox[b] << 16
	y |= sbox[a] << 24
	u := ((y << 13) | (y >> 19)) & 0xFFFFFFFF
	v := ((y << 23) | (y >> 9)) & 0xFFFFFFFF
	return y ^ u ^ v
}

func keySchedule(initialKey []uint32) []uint32 {
	// removed from preview
}

func subDwordRot(x uint32) int {
	// removed from preview
}

func obfuscDword(inval uint32) uint32 {
	bss := new(bytes.Buffer)
	_ = binary.Write(bss, binary.LittleEndian, inval)

	bs := bss.Bytes()

	magic := 0x4E67C6A7
	out := weh(uint32(magic), uint32(bs[0])) ^ uint32(magic)
	out = weh(out, uint32(bs[1])) ^ uint32(out)
	out = weh(out, uint32(bs[2])) ^ uint32(out)
	out = weh(out, uint32(bs[3])) ^ uint32(out)

	return out
}

func mixRowApplyKey(bs []byte, rks []uint32) []uint32 {
	// removed from preview
}

func reverseByteBits(x uint32) uint32 {
	a := ((x >> 1) & 0x55555555) | ((x << 1) & 0xAAAAAAAA)
	b := ((a >> 2) & 0x33333333) | ((a << 2) & 0xCCCCCCCC)
	res := uint32((b>>4)&0xF0F0F0F) | uint32((b<<4)&0xF0F0F0F0)
	return res
}

func mixNibbles(x uint32) uint32 {
	return ((((((x & 0xFF) << 4) + ((x >> 8) & 0xFF)) << 4) + ((x >> 16) & 0xFF)) << 4) + ((x >> 24) & 0xFF)
}

func weh(a uint32, b uint32) uint32 {
	return ((a << 5) + (a >> 2) + b) & 0xFFFFFFFF
}

func rotMess(x uint32) uint32 {

	bss := new(bytes.Buffer)
	_ = binary.Write(bss, binary.LittleEndian, x)

	bs := bss.Bytes()
	out := uint32(bs[0]) << uint32(11)
	out = out | uint32(bs[1])
	out ^= uint32(bs[0]) >> 5
	out ^= 0xFFFFFFFF
	out ^= uint32(bs[0])
	a := (out << 7) & 0xFFFFFFFF
	a ^= uint32(bs[2])
	out ^= (out >> 3) ^ a

	a = (out << 11) & 0xFFFFFFFF
	a ^= uint32(bs[3])
	out ^= ((out >> 5) ^ a) ^ 0xFFFFFFFF
	return out & 0x7FFFFFFF
}

// Endian swap
func reverseBytes(x uint32) uint32 {
	var y uint32
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, x)
	if err != nil {
		fmt.Println("Error writing bytes: ", err)
	}
	b := bytes.NewReader(buf.Bytes())
	_ = binary.Read(b, binary.BigEndian, &y)
	return y
}

func EncryptRB(srcData []byte) string {
	// removed from preview
}

func DecryptRB (encData2 []byte) string {
	// removed from preview
}
