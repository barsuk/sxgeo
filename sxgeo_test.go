package sxgeo

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
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

func TestGetCityFull2(t *testing.T) {
	testCases := []struct {
		ip   string
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

func TestGetCityFull(t *testing.T) {
	SetEndian(LITTLE)
	path := os.Getenv("HOME") + "/Downloads/SxGeoCity.dat"
	_, err := ReadDBToMemory(path)
	if err != nil {
		t.Fatalf("%s %v", path, err)
	}

	testCases := []struct {
		ip string
	}{
		{"31.174.87.24"},
		//{"31.174.87.2",},
		//{"31.174.87.224",},
		//{"198.16.66.100",},
		//{"31.173.87.247",},
		//{"91.193.178.99",},
		//{"178.140.236.47",},
	}
	for _, tc := range testCases {
		fmt.Printf("%s\n", tc.ip)
		c, err := GetCityFull(tc.ip)
		if err != nil {
			t.Fatalf("%v", err)
		}
		fmt.Printf("%+v\n", c)
		enc, err := json.Marshal(c)
		if err != nil {
			fmt.Printf("error: %v", err)
			os.Exit(1)
		}

		fmt.Printf("%s\n", enc)
		os.Exit(0)
	}
}
