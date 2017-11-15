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

	"github.com/alexbeltran/gobacnet"
	"github.com/spf13/cobra"
)

// Flags
var deviceID int

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

func readProp(cmd *cobra.Command, args []string) {
	fmt.Println("readprop called")
	c, err := gobacnet.NewClient(Interface, Port)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// We need the actual address of the device first.
	resp, err := c.WhoIs(deviceID-1, deviceID+1)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)

	//	rp := types.ReadPropertyData{
	//		Object: types.Object{
	//			ID: types.ObjectID{
	//				Type:     0,
	//				Instance: 1,
	//			},
	//			Properties: []types.Property{
	//				types.Property{
	//					Type:       85, // Present value
	//					ArrayIndex: 0xFFFFFFFF,
	//				},
	//			},
	//		},
	//	}
	//
	//	dest := &types.Address{}
	//	c.ReadProperty(dest, rp)
}
func init() {
	RootCmd.AddCommand(readpropCmd)
	readpropCmd.Flags().IntVarP(&deviceID, "device", "d", 1234, "device id")
}
