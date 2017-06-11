package encoding

import "testing"

func TestNPDU(t *testing.T) {
	n := encodeNPDU(false, Normal)
	_, err := EncodePDU(&n, &BacnetAddress{}, &BacnetAddress{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestReadProperty(t *testing.T) {

}

func TestSegsApduEncode(t *testing.T) {
	// Test is structured as parameter 1, parameter 2, output
	tests := [][]int{
		[]int{0, 1, 0},
		[]int{64, 60, 0x61},
		[]int{80, 205, 0x72},
		[]int{80, 405, 0x73},
		[]int{80, 1005, 0x74},
		[]int{3, 1035, 0x15},
		[]int{9, 1035, 0x35},
	}

	for _, test := range tests {
		d := int(encodeMaxSegsMaxApdu(test[0], test[1]))
		if d != test[2] {
			t.Fatalf("Input was Segments %d and Apdu %d: Expected %x got %x", test[0], test[1], test[2], d)
		}
	}
}
