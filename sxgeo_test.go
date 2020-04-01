package sxgeo

import (
	"encoding/binary"
	"testing"
)

func TestSetEndian(t *testing.T) {
	SetEndian(BIG)
	if hbo != binary.BigEndian {
		t.Fatalf("endian not set properly")
	}
	SetEndian(LITTLE)
	if hbo != binary.LittleEndian {
		t.Fatalf("endian not set properly")
	}
}

func TestGetCityFull(t *testing.T) {
	testCases := []struct {
        ip  string
        want string
    }{
        {"224.0.0.0", "IP is loopback or multicast or unspecified"},
		{"127.0.0.1", "IP is loopback or multicast or unspecified"},
	}
    for _, tc := range testCases {
		_, err := GetCityFull(tc.ip)
		if err == nil {
			t.Fatalf("loopback ip should be detected")
		}
    }
}
