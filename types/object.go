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

import "fmt"

type ObjectType uint16

const (
	AnalogInput       ObjectType = 0
	AnalogOutput      ObjectType = 1
	AnalogValue       ObjectType = 2
	BinaryInput       ObjectType = 3
	BinaryOutput      ObjectType = 4
	BinaryValue       ObjectType = 5
	DeviceType        ObjectType = 8
	File              ObjectType = 10
	MultiStateInput   ObjectType = 13
	NotificationClass ObjectType = 15
	MultiStateValue   ObjectType = 19
	TrendLog          ObjectType = 20
	CharacterString   ObjectType = 40
)

const (
	AnalogInputStr       = "Analog Input"
	AnalogOutputStr      = "Analog Output"
	AnalogValueStr       = "Analog Value"
	BinaryInputStr       = "Binary Input"
	BinaryOutputStr      = "Binary Output"
	BinaryValueStr       = "Binary Value"
	DeviceTypeStr        = "Device"
	FileStr              = "File"
	NotificationClassStr = "Notification Class"
	MultiStateValueStr   = "Multi-State Value"
	MultiStateInputStr   = "Multi-State Input"
	TrendLogStr          = "Trend Log"
	CharacterStringStr   = "Character String"
)

var objTypeMap = map[ObjectType]string{
	AnalogInput:       AnalogInputStr,
	AnalogOutput:      AnalogOutputStr,
	AnalogValue:       AnalogValueStr,
	BinaryInput:       BinaryInputStr,
	BinaryOutput:      BinaryOutputStr,
	BinaryValue:       BinaryValueStr,
	DeviceType:        DeviceTypeStr,
	File:              FileStr,
	NotificationClass: NotificationClassStr,
	MultiStateValue:   MultiStateValueStr,
	MultiStateInput:   MultiStateInputStr,
	TrendLog:          TrendLogStr,
	CharacterString:   CharacterStringStr,
}

var objStrTypeMap = map[string]ObjectType{
	AnalogInputStr:       AnalogInput,
	AnalogOutputStr:      AnalogOutput,
	AnalogValueStr:       AnalogValue,
	BinaryInputStr:       BinaryInput,
	BinaryOutputStr:      BinaryOutput,
	BinaryValueStr:       BinaryValue,
	DeviceTypeStr:        DeviceType,
	FileStr:              File,
	NotificationClassStr: NotificationClass,
	MultiStateValueStr:   MultiStateValue,
	TrendLogStr:          TrendLog,
	CharacterStringStr:   CharacterString,
}

func GetType(s string) ObjectType {
	t, ok := objStrTypeMap[s]
	if !ok {
		return 0
	}
	return t
}

func (t ObjectType) String() string {
	s, ok := objTypeMap[t]
	if !ok {
		return fmt.Sprintf("Unknown (%d)", t)
	}
	return fmt.Sprintf("%s", s)
}

type ObjectInstance uint32

type ObjectID struct {
	Type     ObjectType
	Instance ObjectInstance
}

// String returns a pretty print of the ObjectID structure
func (id ObjectID) String() string {
	return fmt.Sprintf("Instance: %d Type: %s", id.Instance, id.Type.String())
}

type Object struct {
	Name        string
	Description string
	ID          ObjectID
	Properties  []Property `json:",omitempty"`
}
