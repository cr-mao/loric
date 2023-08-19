package netlib

import (
	"testing"
)

func TestParseAddr(t *testing.T) {
	listenAddr, exposeAddr, err := ParseAddr(":0")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(listenAddr, exposeAddr)
}

func TestParseAddr2(t *testing.T) {
	listenAddr, exposeAddr, err := ParseAddr("0.0.0.0:12345")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(listenAddr, exposeAddr)
}
