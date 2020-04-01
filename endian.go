package sxgeo

import (
	"encoding/binary"
	"unsafe"
)

// https://stackoverflow.com/a/53286786
func DetectEndian() (binary.ByteOrder, error) {
	buf := [2]byte{}
	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xABCD)

	switch buf {
	case [2]byte{0xCD, 0xAB}:
		print("Little\n")
		return binary.LittleEndian, nil
	case [2]byte{0xAB, 0xCD}:
		print("Big\n")
		return binary.BigEndian, nil
	default:
		panic("Could not determine native endianness.")
	}
}
