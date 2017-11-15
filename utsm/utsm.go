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
	"sync"
	"time"
)

const (
	defaultOverallTimeout = time.Duration(10) * time.Second

	// defaultSubTimeout is how long in between publish packages to a subscriber
	// before we timeout waiting for additional data.
	defaultSubTimeout = time.Duration(1) * time.Second
)

// Manager handles subscriptions and publications. Each manager is thread-safe
type Manager struct {
	subs              []*subscriber
	mutex             *sync.Mutex
	subTimeout        time.Duration
	subOverallTimeout time.Duration
}

// NewManager initializes a manager's internals. Do not allocate a struct of the
// manager directly.
func NewManager(options ...ManagerOption) *Manager {
	m := &Manager{
		subTimeout:        defaultSubTimeout,
		subOverallTimeout: defaultOverallTimeout,
		mutex:             &sync.Mutex{},
	}
	for _, op := range options {
		op(m)
	}
	return m
}

// ManagerOption are function passed to NewManager to configure the manager
type ManagerOption func(m *Manager)

// DefaultSubscriberTimeout option sets a a timeout period when we have not
// received any packages to a subscriber for the timeout period
func DefaultSubscriberTimeout(timeout time.Duration) ManagerOption {
	return func(m *Manager) {
		m.mutex.Lock()
		m.subOverallTimeout = timeout
		m.mutex.Unlock()
	}
}

// DefaultSubscriberTimeout option sets a a timeout period when we have not
// received any packages to a subscriber for the timeout period
func DefaultSubscriberLastReceivedTimeout(timeout time.Duration) ManagerOption {
	return func(m *Manager) {
		m.mutex.Lock()
		m.subTimeout = timeout
		m.mutex.Unlock()
	}
}

func (m *Manager) Publish(id int, data interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, s := range m.subs {
		if id >= s.start && id <= s.end {
			s.mutex.Lock()
			s.lastReceived = time.Now()
			s.data <- data
			s.mutex.Unlock()
		}
	}
}

func (m *Manager) newSubscriber(start int, end int, options []SubscriberOption) *subscriber {
	s := &subscriber{
		start:        start,
		end:          end,
		lastReceived: time.Now(),
		data:         make(chan interface{}, 1),
		mutex:        &sync.Mutex{},
	}
	m.mutex.Lock()
	m.subs = append(m.subs, s)

	s.mutex.Lock()
	s.timeout = m.subOverallTimeout
	s.lastReceivedTimeout = m.subTimeout
	s.mutex.Unlock()

	m.mutex.Unlock()

	for _, opt := range options {
		opt(s)
	}

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
