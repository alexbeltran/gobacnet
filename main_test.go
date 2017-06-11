package gobacnet

import (
	"fmt"
	"testing"
)

// TestMain are general test
func TestMain(t *testing.T) {
	var err error
	_, err = NewClient("wlp1s0")
	if err != nil {
		t.Fatal(err)
	}

	_, err = NewClient("pizzainterfacenotreal")
	if err == nil {
		t.Fatal("Successfully passed a false interface.")
	}
}

func TestGetBroadcast(t *testing.T) {
	failTest := func(addr string) {
		_, err := getBroadcast(addr)
		if err == nil {
			t.Fatalf("%s is not a valid parameter, but it did not gracefully crash", addr)
		}
	}

	failTest("frog")
	failTest("frog/dog")
	failTest("frog/24")
	failTest("16.18.dog/32")

	s, err := getBroadcast("192.168.23.1/24")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(s)
	// Output:
	// 192.168.23.255
}

func TestWhoIs(t *testing.T) {
	c, err := NewClient("wlp1s0")
	if err != nil {
		t.Fatal(err)
	}
	c.sendRequest()
}
