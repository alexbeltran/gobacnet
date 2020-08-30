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

import "encoding/json"

type ObjectMap map[ObjectType]map[ObjectInstance]Object

// Len returns the total number of entries within the object map.
func (o ObjectMap) Len() int {
	counter := 0
	for _, t := range o {
		for _ = range t {
			counter++
		}

	}
	return counter
}

func (om ObjectMap) MarshalJSON() ([]byte, error) {
	m := make(map[string]map[ObjectInstance]Object)
	for typ, sub := range om {
		key := typ.String()
		if m[key] == nil {
			m[key] = make(map[ObjectInstance]Object)
		}
		for inst, obj := range sub {
			m[key][inst] = obj
		}
	}
	return json.Marshal(m)
}

func (om ObjectMap) UnmarshalJSON(data []byte) error {
	m := make(map[string]map[ObjectInstance]Object, 0)
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	for t, sub := range m {
		key := GetType(t)
		if om[key] == nil {
			om[key] = make(map[ObjectInstance]Object)
		}
		for inst, obj := range sub {
			om[key][inst] = obj
		}
	}
	return nil
}
