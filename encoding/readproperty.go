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

package encoding

import (
	"fmt"

	bactype "github.com/alexbeltran/gobacnet/types"
)

func (e *Encoder) readPropertyHeader(tagPos uint8, data bactype.ReadPropertyData) (uint8, error) {
	// Validate data first
	if err := isValidObjectType(data.Object.ID.Type); err != nil {
		return 0, err
	}
	if err := isValidPropertyType(data.Object.Properties[0].Type); err != nil {
		return 0, err
	}

	// Tag - Object Type and Instance
	e.contextObjectID(tagPos, data.Object.ID.Type, data.Object.ID.Instance)
	tagPos++

	// Get first property
	prop := data.Object.Properties[0]
	e.contextEnumerated(tagPos, prop.Type)
	tagPos++

	// Optional Tag - Array Index
	if prop.ArrayIndex != ArrayAll {
		e.contextUnsigned(tagPos, prop.ArrayIndex)
	}
	tagPos++
	return tagPos, nil
}

// ReadProperty is a service request to read a property that is passed.
func (e *Encoder) ReadProperty(invokeID uint8, data bactype.ReadPropertyData) error {
	// PDU Type
	a := bactype.APDU{
		DataType:         bactype.ConfirmedServiceRequest,
		Service:          bactype.ServiceConfirmedReadProperty,
		MaxSegs:          0,
		MaxApdu:          MaxAPDU,
		InvokeId:         invokeID,
		SegmentedMessage: false,
	}
	e.APDU(a)
	e.readPropertyHeader(initialTagPos, data)
	return e.Error()
}

// ReadPropertyAck is the response made to a ReadProperty service request.
func (e *Encoder) ReadPropertyAck(invokeID uint8, data bactype.ReadPropertyData) error {
	if len(data.Object.Properties) != 1 {
		return fmt.Errorf("Property length length must be 1 not %d", len(data.Object.Properties))
	}

	// PDU Type
	a := bactype.APDU{
		DataType: bactype.ComplexAck,
		Service:  bactype.ServiceConfirmedReadProperty,
		MaxSegs:  0,
		MaxApdu:  MaxAPDU,
		InvokeId: invokeID,
	}
	e.APDU(a)

	tagID, err := e.readPropertyHeader(initialTagPos, data)
	if err != nil {
		return err
	}

	e.openingTag(tagID)
	tagID++
	prop := data.Object.Properties[0]
	e.AppData(prop.Data)
	e.closingTag(tagID)
	return e.Error()
}

func (d *Decoder) ReadProperty(data *bactype.ReadPropertyData) error {
	// Must have at least 7 bytes
	if d.buff.Len() < 7 {
		return fmt.Errorf("Missing parameters")
	}

	// Tag 0: Object ID
	tag, meta := d.tagNumber()

	var expectedTag uint8
	if tag != expectedTag {
		return &ErrorIncorrectTag{expectedTag, tag}
	}
	expectedTag++

	var objectType bactype.ObjectType
	var instance bactype.ObjectInstance
	if !meta.isContextSpecific() {
		return fmt.Errorf("Tag %d should be context specific. %x", tag, meta)
	}
	objectType, instance = d.objectId()
	data.Object.ID.Type = objectType
	data.Object.ID.Instance = instance

	// Tag 1: Property ID
	tag, meta = d.tagNumber()
	if tag != expectedTag {
		return &ErrorIncorrectTag{expectedTag, tag}
	}
	expectedTag++

	lenValue := d.value(meta)

	var prop bactype.Property
	prop.Type = d.enumerated(int(lenValue))

	if d.len() != 0 {
		tag, meta = d.tagNumber()
	}

	// Check to see if we still have bytes to read.
	if d.buff.Len() != 0 || tag >= 2 {
		// If we do then that means we are reading the optional argument,
		// arra length

		// Tag 2: Array Length (OPTIONAL)
		var lenValue uint32
		lenValue = d.value(meta)

		var openTag uint8
		// I tried to not use magic numbers but it doesn't look like it can be avoid
		// If the attag we receive is a tag of 2 then set the value
		if tag == 2 {
			prop.ArrayIndex = d.unsigned(int(lenValue))
			if d.len() > 0 {
				openTag, meta = d.tagNumber()
			}
		} else {
			openTag = tag
			prop.ArrayIndex = ArrayAll
		}

		if openTag == 3 {
			var err error
			// We subtract one to ignore the closing tag.
			datalist := make([]interface{}, 0)

			// There is a closing tag of size 1 byte that we ignore which is why we are
			// looping until the length is greater than 1
			for i := 0; d.buff.Len() > 1; i++ {
				data, err := d.AppData()
				if err != nil {
					d.err = err
					return err
				}
				datalist = append(datalist, data)
			}
			prop.Data = datalist

			// If we only have one value in the list, lets just return that value
			if len(datalist) == 1 {
				prop.Data = datalist[0]
			}
			if err != nil {
				d.err = err
				return err
			}
		}
	} else {
		prop.ArrayIndex = ArrayAll
	}

	// We now assemble all the values that we have read above
	data.Object.ID.Instance = instance
	data.Object.ID.Type = objectType
	data.Object.Properties = []bactype.Property{prop}

	return d.Error()
}
