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
	"github.com/alexbeltran/gobacnet/datalink"
	"os"
	"sync"
	"time"

	"github.com/alexbeltran/gobacnet"
	"github.com/alexbeltran/gobacnet/types"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var scanSize uint32
var printStdout bool
var verbose bool
var concurrency int
var output string

var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "discover finds all devices on the network saves results",
	Long:  `discover finds all devices on the network saves results`,

	Run: discover,
}

func save(outfile string, stdout bool, results interface{}) error {
	var file *os.File
	var err error
	if printStdout {
		file = os.Stdout
	} else {
		file, err = os.Create(outfile)

		if err != nil {
			return err
		}
		defer file.Close()
	}
	enc := json.NewEncoder(file)
	enc.SetIndent("", "   ")
	return enc.Encode(results)
}

func discover(cmd *cobra.Command, args []string) {
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{}
	log.Out = os.Stdout
	log.SetLevel(logrus.DebugLevel)

	dataLink, err := datalink.NewUDPDataLink(Interface, Port)
	if err != nil {
		log.Fatal(err)
	}
	c := gobacnet.NewClient(dataLink, 0)
	defer c.Close()
	go c.Run()

	log.Printf("Discovering on interface %s and port %d", Interface, Port)
	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(concurrency)
	scan := make(chan []types.Device, concurrency)
	merge := make(chan types.Device, concurrency)

	// Further discovers new points found in who is
	for i := 0; i < concurrency; i++ {
		go func() {
			for devs := range scan {
				for _, d := range devs {
					log.Infof("Found device: %d", d.ID.Instance)
					dev, err := c.Objects(d)

					if err != nil {
						log.Error(err)
						continue
					}
					merge <- dev
				}
			}
			wg.Done()
		}()

	}

	// combine results
	var results []types.Device
	repeats := make(map[types.ObjectInstance]struct{})
	counter := 0
	total := 0
	go func() {
		for dev := range merge {
			if _, ok := repeats[dev.ID.Instance]; ok {
				log.Errorf("Receive repeated device %d", dev.ID.Instance)
				continue
			}
			log.Infof("Merged: %d", dev.ID.Instance)
			repeats[dev.ID.Instance] = struct{}{}
			if len(dev.Objects) > 0 {
				counter++
			}
			total++
			results = append(results, dev)
		}
	}()

	// Initiates who is
	var startRange, endRange, i int
	incr := int(scanSize)

	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	for i = 0; i < types.MaxInstance/int(scanSize); i++ {
		startRange = i * incr
		endRange = min((i+1)*incr-1, types.MaxInstance)
		log.Infof("Scanning %d to %d", startRange, endRange)
		scanned, err := c.WhoIs(startRange, endRange)
		if err != nil {
			log.Error(err)
			continue
		}
		scan <- scanned
	}
	close(scan)
	wg.Wait()
	close(merge)

	err = save(output, printStdout, results)
	if err != nil {
		log.Errorf("unable to save document: %v", err)
	}
	delta := time.Now().Sub(start)
	log.Infof("Discovery completed in %s", delta)
	if !printStdout {
		log.Infof("Results saved in %s", output)
	}
	log.Infof("%d/%d has values", counter, total)
}

func init() {
	scanSizeDescription := `scan size limits
 the number of devices that are being read at once`

	RootCmd.AddCommand(discoverCmd)
	discoverCmd.Flags().Uint32VarP(&scanSize, "size", "s", 1000, scanSizeDescription)
	discoverCmd.Flags().BoolVar(&printStdout, "stdout", false, "Print to stdout")
	discoverCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Print to additional debugging information")
	discoverCmd.Flags().IntVarP(&concurrency, "concurency", "c", 5, `Number of
  concurrent threads used for scanning the network. A higher number of
  concurrent workers can result in an oversaturate network but will result in
  a faster scan. Concurrency must be greater then 2.`)
	discoverCmd.Flags().StringVarP(&output, "output", "o", "out.json", "Save data to output filename. This field is ignored if stdout is true")
}
