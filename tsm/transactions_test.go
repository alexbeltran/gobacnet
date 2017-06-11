package tsm

import "testing"

func TestTSM(t *testing.T) {
	size := 3
	tsm := New(size)
	var err error
	for i := 0; i < size; i++ {
		_, err = tsm.GetFree()
		if err != nil {
			t.Fatal(err)
		}
	}

	// The buffer should be full at this point.
	_, err = tsm.GetFree()
	if err == nil {
		t.Fatal("Buffer was full but an id was given ")
	}
}
