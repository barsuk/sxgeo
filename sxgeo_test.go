package sxgeo

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
)

var path string

func TestMain(m *testing.M) {
	path = os.Getenv("SXGEODAT") + "/SxGeoCity.dat"
	m.Run()
}

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
	_, err := ReadDBToMemory(path)
	if err != nil {
		t.Fatalf("%s %v", path, err)
	}

	// более-менее валидные адреса можно взять с https://www.4it.me/getlistip?cityid=5138
	tcs := []string{
		"188.255.70.88",
		"1.8.8.8",
		"4.8.8.8",
		"8.9.8.8",
		"192.8.8.8",
		"128.8.8.8",

		"132.95.44.0",

		"31.174.87.24",
		"31.174.87.2",
		"31.174.87.224",
		"198.16.66.100",
		"31.173.87.247",
		"91.193.178.99",
		"178.140.236.47",
		"5.22.153.4",
		"37.49.192.5",
		"2.60.57.9",
		"37.112.130.8",
		"95.107.16.10",
		"24.141.149.0",
	}

	for _, ip := range tcs {
		fmt.Printf("%s\n", ip)
		c, err := GetCityFull(ip)
		if err != nil {
			t.Fatalf("%v", err)
		}
		enc, err := json.Marshal(c)
		if err != nil {
			t.Fatalf("%v", err)
		}

		fmt.Printf("%s\n", enc)
		//os.Exit(0)
	}
}

func TestToMemory(t *testing.T) {
	SetEndian(LITTLE)
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("cannot open DB file: %v", err)
	}
	defer f.Close()

	dbBytes, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("cannot slurp file to slice of bytes %v", err)
	}

	buf := bytes.NewReader(dbBytes)

	_, err = ToMemory(buf)
	if err != nil {
		t.Fatalf("%s %v", path, err)
	}

	// более-менее валидные адреса можно взять с https://www.4it.me/getlistip?cityid=5138
	tcs := []string{
		"188.255.70.88",
		"5.22.153.4",
		"37.49.192.5",
		"2.60.57.9",
		"37.112.130.8",
		"95.107.16.10",
		"24.141.149.0",
	}

	for _, ip := range tcs {
		fmt.Printf("%s\n", ip)
		_, err := GetCityFull(ip)
		if err != nil {
			t.Fatalf("%v", err)
		}
		//enc, err := json.Marshal(c)
		//if err != nil {
		//	t.Fatalf("%v", err)
		//}

		//fmt.Printf("%s\n", enc)
		//os.Exit(0)
	}
}
