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

package gobacnet

import (
	"encoding/json"
	"fmt"
	"github.com/alexbeltran/gobacnet/datalink"
	"github.com/alexbeltran/gobacnet/encoding"
	"log"
	"testing"

	"github.com/alexbeltran/gobacnet/types"
)

const interfaceName = "eth0"
const testServer = 1234

// TestMain are general test
func TestUdpDataLink(t *testing.T) {
	dataLink, err := datalink.NewUDPDataLink(interfaceName, 0)
	if err != nil {
		t.Fatal(err)
	}
	c := NewClient(dataLink, 0)
	c.Close()

	_, err = datalink.NewUDPDataLink("pizzainterfacenotreal", 0)
	if err == nil {
		t.Fatal("Successfully passed a false interface.")
	}
}

func TestMac(t *testing.T) {
	var mac []byte
	json.Unmarshal([]byte("\"ChQAzLrA\""), &mac)
	l := len(mac)
	p := uint16(mac[l-1])<<8 | uint16(mac[l-1])
	log.Printf("%d", p)
}

func TestServices(t *testing.T) {
	dataLink, err := datalink.NewUDPDataLink(interfaceName, 0)
	if err != nil {
		t.Fatal(err)
	}
	c := NewClient(dataLink, 0)
	defer c.Close()

	t.Run("Read Property", func(t *testing.T) {
		testReadPropertyService(c, t)
	})

	t.Run("Who Is", func(t *testing.T) {
		testWhoIs(c, t)
	})

	t.Run("WriteProperty", func(t *testing.T) {
		testWritePropertyService(c, t)
	})

}

func testReadPropertyService(c Client, t *testing.T) {
	dev, err := c.WhoIs(testServer, testServer)
	read := types.PropertyData{
		Object: types.Object{
			ID: types.ObjectID{
				Type:     types.AnalogValue,
				Instance: 1,
			},
			Properties: []types.Property{
				types.Property{
					Type:       types.PropObjectName, // Present value
					ArrayIndex: ArrayAll,
				},
			},
		},
	}
	if len(dev) == 0 {
		t.Fatalf("Unable to find device id %d", testServer)
	}

	resp, err := c.ReadProperty(dev[0], read)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Response: %v", resp.Object.Properties[0].Data)
}

func testWhoIs(c Client, t *testing.T) {
	dev, err := c.WhoIs(testServer-1, testServer+1)
	if err != nil {
		t.Fatal(err)
	}
	if len(dev) == 0 {
		t.Fatalf("Unable to find device id %d", testServer)
	}
}

// This test will first cconver the name of an analogue sensor to a different
// value, read the property to make sure the name was changed, revert back, and
// ensure that the revert was successful
func testWritePropertyService(c Client, t *testing.T) {
	const targetName = "Hotdog"
	dev, err := c.WhoIs(testServer, testServer)
	wp := types.PropertyData{
		Object: types.Object{
			ID: types.ObjectID{
				Type:     types.AnalogValue,
				Instance: 1,
			},
			Properties: []types.Property{
				types.Property{
					Type:       types.PropObjectName, // Present value
					ArrayIndex: ArrayAll,
					Priority:   types.Normal,
				},
			},
		},
	}

	if len(dev) == 0 {
		t.Fatalf("Unable to find device id %d", testServer)
	}
	resp, err := c.ReadProperty(dev[0], wp)
	if err != nil {
		t.Fatal(err)
	}
	// Store the original response since we plan to put it back in after
	org := resp.Object.Properties[0].Data
	t.Logf("original name is: %d", org)

	wp.Object.Properties[0].Data = targetName
	err = c.WriteProperty(dev[0], wp)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = c.ReadProperty(dev[0], wp)
	if err != nil {
		t.Fatal(err)
	}

	d := resp.Object.Properties[0].Data
	s, ok := d.(string)
	if !ok {
		log.Fatalf("unexpected return type %T", d)
	}

	if s != targetName {
		log.Fatalf("write to name %s did not successed, name was %s", targetName, s)
	}

	// Revert Changes
	wp.Object.Properties[0].Data = org
	err = c.WriteProperty(dev[0], wp)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = c.ReadProperty(dev[0], wp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Object.Properties[0].Data != org {
		t.Fatalf("unable to revert name back to original value %v: name is %v", org, resp.Object.Properties[0].Data)
	}
}

func TestDeviceClient(t *testing.T) {
	dataLink, err := datalink.NewUDPDataLink("本地连接", 47809)
	if err != nil {
		fmt.Println(err)
		return
	}
	c := NewClient(dataLink, 0)
	go c.Run()

	devs, err := c.WhoIs(-1, -1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", devs)
	//	c.Objects(devs[0])

	prop, err := c.ReadProperty(
		devs[0],
		types.PropertyData{
			Object: types.Object{
				ID: types.ObjectID{
					Type:     types.AnalogInput,
					Instance: 0,
				},
				Properties: []types.Property{{
					Type:       85,
					ArrayIndex: encoding.ArrayAll,
				}},
			},
			ErrorClass: 0,
			ErrorCode:  0,
		})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(prop.Object.Properties)

	props, err := c.ReadMultiProperty(devs[0], types.MultiplePropertyData{Objects: []types.Object{
		{
			ID: types.ObjectID{
				Type:     types.AnalogInput,
				Instance: 0,
			},
			Properties: []types.Property{
				{
					Type:       8,
					ArrayIndex: encoding.ArrayAll,
				},
				/*	{
					Type:       85,
					ArrayIndex: encoding.ArrayAll,
				},*/
			},
		},
	}})

	fmt.Println(props)
	if err != nil {
		fmt.Println(err)
		return
	}
}
