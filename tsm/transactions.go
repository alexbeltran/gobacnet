package tsm

import "fmt"

const freeID = 0
const MaxTransaction = 255
const InvalidID = 0

const (
	idle = iota
)

type state struct {
	id           int
	state        int
	requestTimer int
}

type TSM struct {
	States []state
	Size   int
	currID int
	count  int
}

// New creates a new transaction manager
func New(size int) *TSM {
	t := TSM{}
	t.Size = size
	t.States = make([]state, size)
	t.currID = 1

	return &t
}

func (t *TSM) incrCursor() {
	t.currID++
	if t.currID == InvalidID {
		t.currID++
	}
}

func (t *TSM) GetFree() (int, error) {
	id, err := t.GetFreeID()
	if err != nil {
		return id, err
	}
	indx, err := t.getFreeIndex()
	if err != nil {
		return id, err
	}

	t.States[indx].id = id
	t.States[indx].state = idle
	t.States[indx].requestTimer = 0 // TODO: apdu_timeout
	t.count = t.count + 1

	return id, nil
}

// GetFreeID returns the first available id. If none is available then MaxTransaction
// is returned
func (t *TSM) GetFreeID() (int, error) {
	if !t.Available() {
		return InvalidID, fmt.Errorf("there are no available ids")
	}
	found := false
	for !found {
		index := t.Find(t.currID)

		// The cursor id is being used, we will skip it
		if index != len(t.States) {
			t.incrCursor()
			continue

			// Cursor is free
		} else {
			id := t.currID
			t.incrCursor()
			return id, nil
		}
	}

	return InvalidID, fmt.Errorf("there are no avialable ids")
}

// getFreeIndex returns the first position in the array that is not being used.
func (t *TSM) getFreeIndex() (int, error) {
	for i, s := range t.States {
		if s.id == InvalidID {
			return i, nil
		}
	}
	return len(t.States), fmt.Errorf("the buffer is full")
}

// Find returns the index where the invoke id has occured.
func (t *TSM) Find(id int) int {
	for i, s := range t.States {
		if s.id == id {
			return i
		}
	}
	return len(t.States)
}

// Avaiable returns true if we can invoke a new id.
func (t *TSM) Available() (status bool) {
	return t.count < len(t.States)
}
