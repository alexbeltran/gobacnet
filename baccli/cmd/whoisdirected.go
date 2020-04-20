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
	"net"
	"os"

	"github.com/alexbeltran/gobacnet"
	"github.com/spf13/cobra"
)

// Flags
var targetIP string
var targetPort int
var outputFilenameWID string

// whoIsCmd represents the whoIs command
var whoIsDirectedCmd = &cobra.Command{
	Use:   "whoIsDirected",
	Short: "BACnet device discovery",
	Long:  `whoIsDirected does a bacnet network discovery directed at a single IP.`,
	Run:   whoIsDirected,
}

func whoIsDirected(cmd *cobra.Command, args []string) {
	c, err := gobacnet.NewClient(Interface, 49999)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	ip := net.ParseIP(targetIP)
	if ip == nil {
		log.Fatal("Unable to parse target IP")
	}

	ids, err := c.WhoIsDirected(ip, targetPort)
	if err != nil {
		log.Fatal(err)
	}

	ioWriter := os.Stdout
	// Check to see if a file was passed to us
	if len(outputFilenameWID) > 0 {
		ioWriter, err = os.Create(outputFilenameWID)
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
	RootCmd.AddCommand(whoIsDirectedCmd)
	whoIsDirectedCmd.Flags().StringVarP(&targetIP, "ip", "x", "", "Target device's IP.")
	whoIsDirectedCmd.Flags().IntVarP(&targetPort, "port", "y", int(0xBAC0), "Target device's port")
	whoIsDirectedCmd.Flags().StringVarP(&outputFilenameWID, "out", "o", "", "Output results into the given filename in json structure.")
}
