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
