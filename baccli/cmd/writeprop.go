// Copyright © 2018 Alex Beltran <alex.e.beltran@gmail.com>
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
	"strconv"

	"github.com/alexbeltran/gobacnet"
	"github.com/alexbeltran/gobacnet/property"
	"github.com/alexbeltran/gobacnet/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// writepropCmd represents the writeprop command
var writepropCmd = &cobra.Command{
	Use:   "writeprop",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: writeProp,
}

// Flags
var targetValue string

func init() {
	// Descriptions are kept separate for legibility purposes.
	propertyTypeDescr := `type of read that will be done. Support both the
	property type as an integer or as a string. e.g. ObjectName or 77 are both
	support. Run --list to see available properties.`
	listPropertiesDescr := `list all string versions of properties that are
	support by property flag`

	RootCmd.AddCommand(writepropCmd)

	// Pass flags to children
	writepropCmd.PersistentFlags().IntVarP(&deviceID, "device", "d", 1234, "device id")
	writepropCmd.Flags().IntVarP(&objectID, "objectID", "o", 1234, "object ID")
	writepropCmd.Flags().IntVarP(&objectType, "objectType", "j", 8, "object type")
	writepropCmd.Flags().StringVarP(&propertyType, "property", "t",
		property.ObjectNameStr, propertyTypeDescr)
	writepropCmd.Flags().StringVarP(&targetValue, "value", "v",
		"", "value that will be set")

	writepropCmd.Flags().Uint32Var(&arrayIndex, "index", gobacnet.ArrayAll, "Which position to return.")
	writepropCmd.PersistentFlags().BoolVarP(&listProperties, "list", "l", false,
		listPropertiesDescr)

}

func writeProp(cmd *cobra.Command, args []string) {
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

	dest := resp[0]

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
				Type:     types.ObjectType(objectType),
				Instance: types.ObjectInstance(objectID),
			},
			Properties: []types.Property{
				types.Property{
					Type:       propInt,
					ArrayIndex: arrayIndex,
				},
			},
		},
	}
	out, err := c.ReadProperty(dest, rp)
	if err != nil {
		if rp.Object.Properties[0].Type == property.ObjectList {
			log.Error("Note: ObjectList reads may need to be broken up into multiple reads due to length. Read index 0 for array length")
		}
		log.Fatal(err)
	}
	if len(out.Object.Properties) == 0 {
		fmt.Println("No value returned")
		return
	}
	data := out.Object.Properties[0].Data
	log.Infof("Current name %v, type %T", out.Object.Properties[0].Data, out.Object.Properties[0].Data)

	if targetValue == "" {
		log.Fatal("nothing was written")
		return
	}

	var wp interface{}
	switch data.(type) {
	case float32:
		var f float64
		f, err = strconv.ParseFloat(targetValue, 32)
		wp = float32(f)
	case float64:
		wp, err = strconv.ParseFloat(targetValue, 32)
	case string:
		wp = targetValue
	default:
		err = fmt.Errorf("unable to handle a type %T", data)
	}
	if err != nil {
		log.Printf("Expects a %T", rp.Object.Properties[0].Data)
	}
	rp.Object.Properties[0].Data = wp
	err = c.WriteProperty(dest, rp, 0)
	if err != nil {
		log.Println(err)
	}
}
