package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

var sbox = []int{
	0x86, 0xBC, 0xFD, 0x94, 0xC9, 0x3F, 0x77, 0x25, 0xC4, 0xE9, 0x9B, 0xD2, 0xA2, 0x59, 0x55,
	0x1E, 0xAB, 0x0C, 0x24, 0x41, 0x4B, 0x4F, 0x9C, 0x2B, 0xE8, 0xE3, 0xE5, 0x03, 0xB9, 0xEC,
	0x21, 0xFC, 0x58, 0x00, 0x30, 0x7B, 0x3A, 0xA8, 0x1C, 0x53, 0x11, 0x6E, 0x9D, 0xB6, 0xF9,
	0x70, 0xB5, 0xE1, 0x43, 0x09, 0xB4, 0x08, 0xBD, 0x3B, 0xA0, 0xFF, 0x80, 0xA9, 0x29, 0x32,
	0xBE, 0x15, 0x65, 0xEA, 0x7C, 0x26, 0xED, 0xDD, 0x9E, 0x0D, 0x5F, 0x4E, 0x37, 0x2E, 0xA5,
	0x0A, 0x01, 0x97, 0x7F, 0x6D, 0x5A, 0x89, 0x62, 0x4C, 0xE7, 0x5D, 0x2D, 0x8C, 0xAF, 0x95,
	0x13, 0x7E, 0xCB, 0x8E, 0x0F, 0xC1, 0xB3, 0xC8, 0x19, 0xF8, 0xC7, 0x8F, 0x04, 0x1A, 0xBF,
	0x52, 0x34, 0x9F, 0x2A, 0x79, 0x56, 0x81, 0xD0, 0x98, 0x91, 0x39, 0x02, 0xB1, 0x38, 0x0B,
	0x12, 0xBA, 0xF6, 0x6F, 0x5E, 0xEB, 0x3D, 0xD7, 0xE6, 0xFE, 0x50, 0xE0, 0x6A, 0x87, 0x85,
	0x1F, 0xDC, 0x64, 0x71, 0xB0, 0x3E, 0xA4, 0x63, 0xF0, 0x68, 0x83, 0x46, 0x93, 0xEE, 0x60,
	0x82, 0xAA, 0x88, 0x10, 0x40, 0x22, 0x66, 0xC0, 0x5C, 0x2C, 0xDB, 0x54, 0xC2, 0x61, 0x7D,
	0x20, 0x1D, 0x8B, 0x84, 0x8A, 0x7A, 0x4A, 0x18, 0x33, 0xF1, 0xFB, 0x06, 0x49, 0x51, 0x74,
	0x73, 0x17, 0xB7, 0xD6, 0x9A, 0xD8, 0xF5, 0x14, 0xCC, 0xB8, 0x0E, 0x4D, 0xDE, 0xBB, 0x42,
	0xF4, 0xFA, 0x2F, 0xDA, 0xAE, 0xA1, 0x90, 0xB2, 0x31, 0x35, 0x05, 0xD9, 0xC5, 0xCD, 0xDF,
	0x44, 0x07, 0x76, 0x6C, 0x27, 0xAD, 0x67, 0x23, 0xD1, 0x92, 0x28, 0xCF, 0x78, 0x75, 0xC3,
	0xAC, 0x69, 0x48, 0xF3, 0xE4, 0x36, 0x72, 0x5B, 0x3C, 0x1B, 0x47, 0xCA, 0xC6, 0x6B, 0xF2,
	0xA6, 0xF7, 0xD5, 0x99, 0x96, 0xD4, 0x45, 0xA3, 0x16, 0xCE, 0xA7, 0x57, 0x8D, 0xE2, 0xEF,
	0xD3,
}

var blowfishSub = []int{
	0x9C30D539,
	0x2AF26013,
	0xC5D1B023,
	0x286085F0,
}

var blowfishP = []int{
	0x452821E6, 0x38D01377, 0xBE5466CF, 0x34E90C6C, 0xC0AC29B7, 0xC97C50DD,
	0x3F84D5B5, 0xB5470917, 0x9216D5D9, 0x8979FB1B, 0xD1310BA6, 0x98DFB5AC,
	0x2FFD72DB, 0xD01ADFB7, 0xB8E1AFED, 0x6A267E96, 0xBA7C9045, 0xF12C7F99,
	0x24A19947, 0xB3916CF7, 0x801F2E2, 0x858EFC16, 0x636920D8, 0x71574E69,
	0xA458FEA3, 0xF4933D7E, 0xD95748F, 0x728EB658, 0x718BCD58, 0x82154AEE,
	0x7B54A41D, 0xC25A59B5,
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
	cs := []uint32{0, 0, 0, 0}
	cs[0] = initialKey[0] ^ uint32(blowfishSub[0])
	cs[1] = initialKey[1] ^ uint32(blowfishSub[1])
	cs[2] = initialKey[2] ^ uint32(blowfishSub[2])
	cs[3] = initialKey[3] ^ uint32(blowfishSub[3])
	state := make([]uint32, 32)
	for i := range state {
		x := cs[0] ^ cs[1] ^ cs[2] ^ cs[3] ^ uint32(blowfishP[i])
		state[i] = uint32(subDword(int(x)))
		cs[0] = cs[1]
		cs[1] = cs[2]
		cs[2] = cs[3]
		cs[3] = uint32(state[i])
	}
	return state
}

func subDwordRot(x uint32) int {
	a := (x >> 24) & 0xFF
	b := (x >> 16) & 0xFF
	c := (x >> 8) & 0xFF
	d := x & 0xFF
	v1 := sbox[a]
	v2 := sbox[d]
	v3 := (sbox[b] << 16) | (v1 << 24)
	v4 := v3 | (sbox[c] << 8)
	out := (((v4 | v2) << 10) | (v3 >> 22)) ^ ((v4 >> 8) | (v2 << 24)) ^ (v4 | v2) ^ (4*(v4|v2) | (v1 >> 6)) ^ (((v4 | v2) << 18) | (v4 >> 14))
	return out & 0xFFFFFFFF
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
	var dws [4]uint32
	bb := bytes.NewReader(bs)
	_ = binary.Read(bb, binary.BigEndian, &dws)

	// dws := bs[0:16]
	state := []uint32{
		dws[0],
		dws[1],
		dws[2],
		dws[3],
	}

	for i := 0; i < 32; i++ {
		t := rks[i] ^ state[3] ^ state[2] ^ state[1]
		out := uint32(subDwordRot(t)) ^ state[0]
		state[0] = state[1]
		state[1] = state[2]
		state[2] = state[3]
		state[3] = out
	}

	return state
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
	// pack as little endian uint32 (should already be as such)
	// Read x as characters
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
	nblocks := len(srcData) / 16
	var outbs []byte
	initialKey := []uint32{
		0x727167A2,
		0xD1E4033C,
		0x47D4044B,
		0xFD850DD2,
	}
	roundKey := keySchedule(initialKey)
	for i := 0; i < nblocks; i++ {
		low := i * 16
		high := (i + 1) * 16
		blockbs := srcData[low:high]
		var dws [4]uint32
		bb := bytes.NewReader(blockbs)
		_ = binary.Read(bb, binary.LittleEndian, &dws)

		// reverse the bits
		var dws2 [4]uint32
		for d := range dws {
			dws2[d] = reverseByteBits(dws[d])
		}

		nb := []uint32{
			dws2[3],
			dws2[2],
			dws2[1],
			dws2[0],
		}

		buf := new(bytes.Buffer)

		for i := range nb {
			err := binary.Write(buf, binary.BigEndian, nb[i])
			if err != nil {
				fmt.Println("Error writing bytes: ", err)
			}
		}

		blockbs = buf.Bytes()
		ia := make([]uint32, len(roundKey))
		for i := range roundKey {
			ia[i] = uint32(roundKey[i])
		}

		state := mixRowApplyKey(blockbs, ia)
		nstate := make([]uint32, len(state))
		for i := range state {
			xx := state[i]
			nstate[i] = reverseBytes(xx)
		}
		state = nstate

		for j := 0; j < 4; j++ {
			tmpstate := make([]uint32, len(state))
			x := reverseByteBits(mixNibbles(state[2]))
			y := reverseByteBits(rotMess(state[1]))
			z := reverseByteBits(obfuscDword(state[0]))
			w := reverseByteBits(state[3])
			out := x ^ y ^ z ^ w
			bss := new(bytes.Buffer)
			_ = binary.Write(bss, binary.LittleEndian, out)
			bs := bss.Bytes()
			outbs = append(outbs, bs...)
			tmpstate[1] = state[0]
			tmpstate[2] = state[1]
			tmpstate[3] = state[2]
			tmpstate[0] = out
			state = tmpstate
		}
		dws3 := make([]uint32, 4)
		bbl := bytes.NewReader(blockbs)
		_ = binary.Read(bbl, binary.BigEndian, &dws3)
		roundKey = keySchedule(dws3)
	}

	rem := len(srcData) % 16
	for i := 0; i < rem; i++ {
		targetIndex := i-1
		targetChar := srcData[:len(srcData)-targetIndex-1]
		ch := targetChar[len(targetChar)-1]
		if rem % 2 == 1 && i == (rem/2) {
			outbs = append(outbs, ch)
			continue
		}
		rev := reverseByteBits(uint32(ch))
		outbs = append(outbs, byte(rev))
	}
	result := fmt.Sprintf("01%02x", outbs)
	return result
}

func DecryptRB (encData2 []byte) string {
	encData, _ := hex.DecodeString(string(encData2))
	encData = encData[1:]
	nblocks := len(encData) / 16
	var outbs []byte
	initialKey := []uint32{
		0x727167A2,
		0xD1E4033C,
		0x47D4044B,
		0xFD850DD2,
	}
	roundKey := keySchedule(initialKey)
	for i := 0; i < nblocks; i++ {
		low := i * 16
		high := (i + 1) * 16
		blockbs := encData[low:high]
		var dws [4]uint32
		bb := bytes.NewReader(blockbs)
		_ = binary.Read(bb, binary.LittleEndian, &dws)
		buf := new(bytes.Buffer)
		_ = binary.Write(buf, binary.BigEndian, dws[3])
		_ = binary.Write(buf, binary.BigEndian, dws[2])
		_ = binary.Write(buf, binary.BigEndian, dws[1])
		_ = binary.Write(buf, binary.BigEndian, dws[0])
		blockbs = buf.Bytes()
		state := []uint32{dws[0],dws[1],dws[2],dws[3]}
		for j := 0; j < 4; j++ {
			tmpstate := make([]uint32, len(state))
			x := reverseByteBits(mixNibbles(state[0]))
			y := reverseByteBits(rotMess(state[1]))
			z := reverseByteBits(obfuscDword(state[2]))
			out := reverseByteBits(x^y^z^state[3])
			tmpstate[1] = state[0]
			tmpstate[2] = state[1]
			tmpstate[3] = state[2]
			tmpstate[0] = out
			state = tmpstate
		}
		ts := make([]uint32, len(state))
		for i, x := range state {
			ts[i] = reverseBytes(x)
		}
		state = ts
		for i := range state {
			blockbs = append(blockbs, byte(state[i]))
		}
		buf2 := new(bytes.Buffer)
		for i := range state {
			err := binary.Write(buf2, binary.BigEndian, state[i])
			if err != nil {
				fmt.Println("Error writing bytes: ", err)
			}
		}
		blockbs = buf2.Bytes()
		s := roundKey
		for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
			s[i], s[j] = s[j], s[i]
		}
		state = mixRowApplyKey(blockbs, s)
		mixed := []uint32{state[3],state[2],state[1],state[0]}
		roundKey = keySchedule(mixed)
		for i := range state {
			state[i] = reverseByteBits(state[i])
		}
		tmpBuffer := new(bytes.Buffer)
		_ = binary.Write(tmpBuffer, binary.LittleEndian, state[0])
		_ = binary.Write(tmpBuffer, binary.LittleEndian, state[1])
		_ = binary.Write(tmpBuffer, binary.LittleEndian, state[2])
		_ = binary.Write(tmpBuffer, binary.LittleEndian, state[3])
		tmpBytes := tmpBuffer.Bytes()
		outbs = append(outbs, tmpBytes...)
	}
	rem := len(encData)%16
	for i := 0; i < rem; i++ {
		targetIndex := i-1
		targetChar := encData[:len(encData)-targetIndex-1]
		ch := targetChar[len(targetChar)-1]
		if rem % 2 == 1 && i == (rem/2) {
			outbs = append(outbs, ch)
			continue
		}
		rev := reverseByteBits(uint32(ch))
		outbs = append(outbs, byte(rev))
	}
	return string(outbs)[4:] // Nuke the first 4 bytes as it's the padding
}
