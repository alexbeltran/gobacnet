// Copyright © 2017 Alex Beltran <alex.e.beltran@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/viper"

	"github.com/alexbeltran/gobacnet"
	"github.com/alexbeltran/gobacnet/property"
	"github.com/alexbeltran/gobacnet/types"
	"github.com/spf13/cobra"
)

// Flags
var (
	deviceID       int
	objectID       int
	objectType     int
	propertyType   string
	listProperties bool
)

// readpropCmd represents the readprop command
var readpropCmd = &cobra.Command{
	Use:   "readprop",
	Short: "Prints out a device's object's property",
	Long: `
 Given a device's object instance and selected property, we print the value
 stored there. There are some autocomplete features to try and minimize the
 amount of arguments that need to be passed, but do take into consideration
 this discovery process may cause longer reads.
	`,
	Run: readProp,
}

func readProp(cmd *cobra.Command, args []string) {
	if listProperties {
		property.PrintAll()
		return
	}

	c, err := gobacnet.NewClient(viper.GetString("interface"), viper.GetInt("port"))
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// We need the actual address of the device first.
	resp, err := c.WhoIs(deviceID, deviceID)
	if err != nil {
		log.Fatal(err)
	}

	if len(resp) == 0 {
		log.Fatal("Device id was not found on the network.")
	}

	dest := &resp[0]

	var propInt uint32
	// Check to see if an int was passed
	if i, err := strconv.Atoi(propertyType); err == nil {
		propInt = uint32(i)
	} else {
		propInt, err = property.Get(propertyType)
	}

	if property.IsDeviceProperty(propInt) {
		objectType = 8
	}

	if err != nil {
		log.Fatal(err)
	}

	rp := types.ReadPropertyData{
		Object: types.Object{
			ID: types.ObjectID{
				Type:     uint16(objectType),
				Instance: uint32(objectID),
			},
			Properties: []types.Property{
				types.Property{
					Type:       propInt,
					ArrayIndex: 0xFFFFFFFF,
				},
			},
		},
	}
	out, err := c.ReadProperty(&dest.Addr, rp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out.Object.Properties[0].Data)
}
func init() {
	// Descriptions are kept separate for legibility purposes.
	propertyTypeDescr := `type of read that will be done. Support both the
	property type as an integer or as a string. e.g. ObjectName or 77 are both
	support. Run --list to see available properties.`
	listPropertiesDescr := `list all string versions of properties that are
	support by property flag`

	RootCmd.AddCommand(readpropCmd)

	// Pass flags to children
	readpropCmd.PersistentFlags().IntVarP(&deviceID, "device", "d", 1234, "device id")
	readpropCmd.Flags().IntVarP(&objectID, "objectID", "o", 1234, "object ID")
	readpropCmd.Flags().IntVarP(&objectType, "objectType", "j", 8, "object type")
	readpropCmd.Flags().StringVarP(&propertyType, "property", "t",
		property.ObjectNameStr, propertyTypeDescr)

	readpropCmd.PersistentFlags().BoolVarP(&listProperties, "list", "l", false,
		listPropertiesDescr)
}
