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
	"strings"

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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: readProp,
}

func longestString(x map[string]uint32) int {
	max := 0
	for k, _ := range x {
		if len(k) > max {
			max = len(k)
		}
	}
	return max
}

const numOfAdditionalSpaces = 15

func printRow(col1, col2 string, maxLen int) {
	spacing := strings.Repeat(" ", maxLen-len(col1)+numOfAdditionalSpaces)
	fmt.Printf("%s%s%s\n", col1, spacing, col2)
}

func printKeys() {
	keys := property.Keys()
	max := longestString(keys)

	printRow("Key", "Int", max)
	fmt.Println(strings.Repeat("-", max+numOfAdditionalSpaces+6))

	for k, id := range keys {
		// Spacing
		printRow(k, fmt.Sprintf("%d", id), max)
	}
}

func readProp(cmd *cobra.Command, args []string) {
	if listProperties {
		printKeys()
		return
	}

	fmt.Println("readprop called")
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
	log.Printf("Getting ID: %d Type %d Prop %d Index %d",
		rp.Object.ID.Instance, rp.Object.ID.Type, rp.Object.Properties[0].Type,
		rp.Object.Properties[0].ArrayIndex)

	out, err := c.ReadProperty(&dest.Addr, rp)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("SUCCESS")
	log.Println(out.Object.Properties[0].Data)
}
func init() {
	RootCmd.AddCommand(readpropCmd)
	readpropCmd.Flags().IntVarP(&deviceID, "device", "d", 1234, "device id")
	readpropCmd.Flags().IntVarP(&objectID, "objectID", "o", 1234, "object ID")
	readpropCmd.Flags().IntVarP(&objectType, "objectType", "j", 8, "object type")
	readpropCmd.Flags().StringVarP(&propertyType, "property", "t",
		property.ObjectNameStr, "type of read that will be done. Support both"+
			"the property type as an integer or as a string. Run --list to see available"+
			"properties.")
	readpropCmd.Flags().BoolVarP(&listProperties, "list", "l", false, "list all properties")
}
