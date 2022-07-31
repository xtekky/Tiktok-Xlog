// Xlog v01
// @contact: t.me/xtekky

package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

var sbox = []int{
	// preview
}

var blowfishSub = []int{
	// preview
}

var blowfishP = []int{
	// preview
}

func subDword(x int) int {
	// preview
}

func keySchedule(initialKey []uint32) []uint32 {
	// preview
}

func subDwordRot(x uint32) int {
	// preview
}

func obfuscDword(inval uint32) uint32 {
	// preview
}

func mixRowApplyKey(bs []byte, rks []uint32) []uint32 {
	// preview
}

func reverseByteBits(x uint32) uint32 {
	// preview
}

func mixNibbles(x uint32) uint32 {
	// preview
}

func weh(a uint32, b uint32) uint32 {
	// preview
}

func rotMess(x uint32) uint32 {
	// preview
}

// Endian swap
func reverseBytes(x uint32) uint32 {
	// preview
}

func EncryptRB(srcData []byte) string {
	// preview
}

func DecryptRB (encData2 []byte) string {
	// preview
}
