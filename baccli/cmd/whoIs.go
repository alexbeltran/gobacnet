// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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
// See the License for the specific language governing permissions and // limitations under the License.

package cmd

import (
	"encoding/json"
	"log"
	"os"

	"github.com/alexbeltran/gobacnet"
	"github.com/spf13/cobra"
)

// Flags
var Interface string
var Port int
var startRange int
var endRange int
var outputFilename string

// whoIsCmd represents the whoIs command
var whoIsCmd = &cobra.Command{
	Use:   "whoIs",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: main,
}

func main(cmd *cobra.Command, args []string) {
	c, err := gobacnet.NewClient(Interface, Port)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	ids, err := c.WhoIs(startRange, endRange)
	if err != nil {
		log.Fatal(err)
	}

	ioWriter := os.Stdout
	// Check to see if a file was passed to us
	if len(outputFilename) > 0 {
		ioWriter, err = os.Create(outputFilename)
		if err != nil {
			log.Fatal(err)
		}
		defer ioWriter.Close()
	}
	// Pretty Print!
	w := json.NewEncoder(ioWriter)
	w.SetIndent("", "    ")
	w.Encode(ids)

}

func init() {
	RootCmd.AddCommand(whoIsCmd)
	whoIsCmd.Flags().StringVarP(&Interface, "interface", "i", "eth0", "Interface e.g. eth0")
	whoIsCmd.Flags().IntVarP(&Port, "port", "p", int(0xBAC0), "Port")
	whoIsCmd.Flags().IntVarP(&startRange, "start", "s", -1, "Start range of discovery")
	whoIsCmd.Flags().IntVarP(&endRange, "end", "e", int(0xBAC0), "End range of discovery")
	whoIsCmd.Flags().StringVarP(&outputFilename, "out", "o", "", "Output results into the given filename in json structure.")
}
