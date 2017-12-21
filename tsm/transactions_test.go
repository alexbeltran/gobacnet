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
	"testing"
	"time"
)

func TestTSM(t *testing.T) {
	size := 3
	tsm := New(size)
	ctx := context.Background()
	var err error
	for i := 0; i < size-1; i++ {
		_, err = tsm.ID(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}

	id, err := tsm.ID(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// The buffer should be full at this point.
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond)
	defer cancel()
	_, err = tsm.ID(ctx)
	if err == nil {
		t.Fatal("Buffer was full but an id was given ")
	}

	// Free an ID
	err = tsm.Put(id)
	if err != nil {
		t.Fatal(err)
	}

	// Now we should be able to get a new id since we free id
	_, err = tsm.ID(context.Background())
	if err != nil {
		t.Fatal(err)
	}

}

func TestDataTransaction(t *testing.T) {
	size := 2
	tsm := New(size)
	ids := make([]int, size)
	var err error

	for i := 0; i < size; i++ {
		ids[i], err = tsm.ID(context.Background())
		if err != nil {
			t.Fatal(err)
		}
	}

	go func() {
		err = tsm.Send(ids[0], "Hello First ID")
		if err != nil {
			t.Error(err)
		}
	}()

	go func() {
		err = tsm.Send(ids[1], "Hello Second ID")
		if err != nil {
			t.Error(err)
		}
	}()

	go func() {
		b, err := tsm.Receive(ids[0], time.Duration(5)*time.Second)
		if err != nil {
			t.Error(err)
		}
		s, ok := b.(string)
		if !ok {
			t.Errorf("type was not preseved")
			return
		}
		t.Log(s)
	}()

	b, err := tsm.Receive(ids[1], time.Duration(5)*time.Second)
	if err != nil {
		t.Error(err)
	}

	s, ok := b.(string)
	if !ok {
		t.Errorf("type was not preseved")
		return
	}
	t.Log(s)
}
