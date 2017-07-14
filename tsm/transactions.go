/*Copyright (C) 2017 Alex Beltran

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to:
The Free Software Foundation, Inc.
59 Temple Place - Suite 330
Boston, MA  02111-1307, USA.

As a special exception, if other files instantiate templates or
use macros or inline functions from this file, or you compile
this file and link it with other works to produce a work based
on this file, this file does not by itself cause the resulting
work to be covered by the GNU General Public License. However
the source code for this file must still be made available in
accordance with section (3) of the GNU General Public License.

This exception does not invalidate any other reasons why a work
based on this file might be covered by the GNU General Public
License.
*/

package tsm

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const freeID = 0

// MaxTransaction is the default maximium number of transactions that can occur
// concurrently
const MaxTransaction = 255
const invalidID = 0

const (
	idle = iota
)

type state struct {
	id           int
	state        int
	requestTimer int
	data         chan []byte
}

// TSM is the transaction state manager. It handles passing data to other
// processes and keeping track of what transactions are currently processed
type TSM struct {
	states []state
	size   int
	currID int
	count  int
	mutex  *sync.Mutex
}

// New creates a new transaction manager
func New(size int) *TSM {
	t := TSM{}
	t.size = size
	t.states = make([]state, size)
	t.currID = 1
	t.mutex = &sync.Mutex{}

	// Initialize the channel pipeline
	for i := range t.states {
		t.states[i].data = make(chan []byte, 1)
	}
	return &t
}

// incrCursor moves the current possible id by one. It handles wrap around
func (t *TSM) incrCursor() {
	t.currID++
	if t.currID == invalidID || t.currID >= MaxTransaction {
		t.currID = invalidID + 1
	}
}

// Send data to invoked id
func (t *TSM) Send(id int, b []byte) error {
	t.mutex.Lock()
	i, err := t.find(id)
	t.mutex.Unlock()
	if err != nil {
		return err
	}
	t.states[i].data <- b

	return nil
}

// Receive attempts to receive a byte array from the invoked id. If a time out
// period has passed then an error is returned
func (t *TSM) Receive(id int, timeout time.Duration) ([]byte, error) {
	t.mutex.Lock()
	i, err := t.find(id)
	t.mutex.Unlock()

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Wait for data
	select {
	case b := <-t.states[i].data:
		return b, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("Receive timed out")
	}

}

// GetFree returns the invoke id that was used to save the state of this connection.
func (t *TSM) GetFree() (int, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	id, err := t.getFreeID()
	if err != nil {
		return id, err
	}
	indx, err := t.getFreeIndex()
	if err != nil {
		return id, err
	}

	t.states[indx].id = id
	t.states[indx].state = idle
	t.states[indx].requestTimer = 0 // TODO: apdu_timeout
	t.count = t.count + 1

	return id, nil
}

// FreeID allows the id to be reused in the transaction manager
func (t *TSM) FreeID(id int) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	i, err := t.find(id)
	if err != nil {
		return err
	}

	t.states[i].id = invalidID
	t.count--
	return nil
}

// GetFreeID returns the first available id. If none is available then MaxTransaction
// is returned
func (t *TSM) getFreeID() (int, error) {
	if !t.available() {
		return invalidID, fmt.Errorf("The TSM is full, there are no available ids")
	}
	found := false
	counter := 0
	for !found && counter < MaxTransaction {
		_, err := t.find(t.currID)

		// The cursor id is being used, we will skip it
		if err == nil {
			t.incrCursor()
			// Cursor is free
		} else {
			id := t.currID
			t.incrCursor()
			return id, nil
		}
		counter++
	}

	return invalidID, fmt.Errorf("there are no available ids")
}

// getFreeIndex returns the first position in the array that is not being used.
func (t *TSM) getFreeIndex() (int, error) {
	for i, s := range t.states {
		if s.id == invalidID {
			return i, nil
		}
	}
	return len(t.states), fmt.Errorf("the buffer is full")
}

// find returns the index where the invoke id has occurred.
func (t *TSM) find(id int) (int, error) {
	for i, s := range t.states {
		if s.id == id {
			return i, nil
		}
	}
	return 0, fmt.Errorf("Unable to find index")
}

// available returns true if we can invoke a new id.
func (t *TSM) available() (status bool) {
	return t.count < len(t.states)
}
