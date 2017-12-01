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
	"encoding/json"
	"log"
	"os"

	"github.com/alexbeltran/gobacnet"
	"github.com/alexbeltran/gobacnet/types"
	"github.com/spf13/cobra"
)

var scanSize uint32
var printStdout bool

// discoverCmd represents the discover command
var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "discover finds all devices on the network saves results",
	Long:  `discover finds all devices on the network saves results`,

	Run: discover,
}

const maxWorkerPool = 20

// gather is a worker who manages getting each devices objects and combining it
// to a completed channel
func gather(client *gobacnet.Client, devChan chan []types.Device, complete chan []types.Device) {
	var devices []types.Device
	var err error
	for devs := range devChan {
		results := make(chan error, len(devs))

		// Query all of the devices
		for i, d := range devs {
			go func(i int, d types.Device) {
				log.Printf("Discover Device: %d", d.ID.Instance)
				devs[i], err = client.Objects(d)
				results <- err
			}(i, d)
		}

		// Aggregate the responses
		for i := 0; i < len(devs); i++ {
			err = <-results
			if err != nil {
				log.Println(err)
			}
		}
		devices = append(devices, devs...)
	}
	complete <- devices
}

func scan(c *gobacnet.Client, startRange, endRange int, dev chan []types.Device) error {
	log.Printf("Scanning %d to %d\n", startRange, endRange)
	subsetDevices, err := c.WhoIs(startRange, endRange)
	if err != nil {
		return err
	}
	dev <- subsetDevices
	return nil
}

func discover(cmd *cobra.Command, args []string) {
	c, err := gobacnet.NewClient(Interface, Port)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	devChan := make(chan []types.Device, maxWorkerPool)
	finish := make(chan []types.Device, 1)
	go gather(c, devChan, finish)
	incr := int(scanSize)
	i := 0

	var startRange, endRange int
	for i = 0; i < types.MaxInstance/int(incr); i++ {
		startRange = i * incr
		endRange = (i+1)*incr - 1
		err = scan(c, startRange, endRange, devChan)
		if err != nil {
			log.Print(err)
		}
	}

	startRange = i * incr
	endRange = types.MaxInstance
	err = scan(c, startRange, endRange, devChan)
	if err != nil {
		log.Print(err)
	}

	close(devChan)
	out := <-finish
	var file *os.File
	if printStdout {
		file = os.Stdout
	} else {
		file, err = os.Create("test.json")

		if err != nil {
			log.Println(err)
		}
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "   ")
	enc.Encode(out)
}

func init() {
	scanSizeDescription := `scan size limits
 the number of devices that are being read at once`

	RootCmd.AddCommand(discoverCmd)
	discoverCmd.Flags().Uint32VarP(&scanSize, "size", "s", 10000, scanSizeDescription)
	discoverCmd.Flags().BoolVar(&printStdout, "stdout", false, "Print to stdout")
}
