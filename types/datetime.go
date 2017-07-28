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

package types

type DayOfWeek int

const (
	None      DayOfWeek = iota
	Monday    DayOfWeek = iota
	Tuesday   DayOfWeek = iota
	Wednesday DayOfWeek = iota
	Thursday  DayOfWeek = iota
	Friday    DayOfWeek = iota
	Saturday  DayOfWeek = iota
	Sunday    DayOfWeek = iota
)

type Date struct {
	Year  int
	Month int
	Day   int
	// Bacnet has an option to only do operations on even or odd months
	EvenMonth      bool
	OddMonth       bool
	EvenDay        bool
	OddDay         bool
	LastDayOfMonth bool
	DayOfWeek      DayOfWeek
}

type Time struct {
	Hour        int
	Minute      int
	Second      int
	Millisecond int
}

// UnspecifiedTime means that this time is triggered through out a period. An
// example of this is 02:FF:FF:FF will trigger all through out 2 am
const UnspecifiedTime = 0xFF
