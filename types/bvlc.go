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

// BACnet Virtual Link Control (BVLC)

// BVLCTypeBacnetIP is the only valid type for the BVLC layer as of 2002.
// Additional types may be added in the future
const BVLCTypeBacnetIP = 0x81

// Bacnet Fuction
type BacFunc byte

// List of possible BACnet functions
const (
	BacFuncResult                          BacFunc = 0
	BacFuncWriteBroadcastDistributionTable BacFunc = 1
	BacFuncBroadcastDistributionTable      BacFunc = 2
	BacFuncBroadcastDistributionTableAck   BacFunc = 3
	BacFuncForwardedNPDU                   BacFunc = 4
	BacFuncUnicast                         BacFunc = 10
	BacFuncBroadcast                       BacFunc = 11
)

type BVLC struct {
	Type     byte
	Function BacFunc

	// Length includes the length of Type, Function, and Length. (4 bytes) It also
	// has the length of the data field after
	Length uint16
	Data   []byte
}
