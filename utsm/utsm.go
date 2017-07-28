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

package utsm

import (
	"bytes"
	"context"
	"log"
	"sync"
	"time"
)

type subscriber struct {
	// Start and End is the range that this object is subscribed to
	Start        int
	End          int
	LastReceived time.Time
	// Data channel is used for data transfer between subscriber and publisher
	Data  chan []byte
	mutex *sync.Mutex
}

type Manager struct {
	subs  []*subscriber
	mutex *sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		mutex: &sync.Mutex{},
	}
}

func (s *subscriber) Deadline() time.Duration {
	s.mutex.Lock()
	// Deadline is x seconds after the last packet we received.
	deadline := s.LastReceived.Add(time.Duration(1) * time.Second).Sub(time.Now())
	s.mutex.Unlock()
	return deadline
}

func (m *Manager) Publish(id int, data []byte) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for i, s := range m.subs {
		if id >= s.Start && id <= s.End {
			log.Printf("%d", i)
			s.mutex.Lock()
			s.Data <- data
			s.LastReceived = time.Now()
			s.mutex.Unlock()
		}
	}
}

func (m *Manager) newSubscriber(start int, end int) *subscriber {
	s := &subscriber{
		Start:        start,
		End:          end,
		LastReceived: time.Now(),
		Data:         make(chan []byte, 1),
		mutex:        &sync.Mutex{},
	}
	m.mutex.Lock()
	m.subs = append(m.subs, s)
	m.mutex.Unlock()
	return s
}

func (m *Manager) removeSubscriber(sub *subscriber) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for i, s := range m.subs {
		if s == sub {
			// https://github.com/golang/go/wiki/SliceTricks
			// Prevents a memory leak that may occur when deleting

			// Shift
			copy(m.subs[i:], m.subs[i+1:])

			// Set last value nil
			m.subs[len(m.subs)-1] = nil

			// Remove last value
			m.subs = m.subs[:len(m.subs)-1]
			return
		}
	}
}

// Subscribe
func (m *Manager) Subscribe(start int, end int, timout time.Duration) ([]byte, error) {
	var buff bytes.Buffer
	s := m.newSubscriber(start, end)
	defer m.removeSubscriber(s)

	ctx, cancel := context.WithTimeout(context.Background(), timout)
	defer cancel()

	for {
		c, can := context.WithTimeout(ctx, s.Deadline())
		defer can()

		select {
		case <-c.Done():
			return buff.Bytes(), nil
		case b := <-s.Data:
			buff.Write(b)
		}
	}
	return buff.Bytes(), nil
}
