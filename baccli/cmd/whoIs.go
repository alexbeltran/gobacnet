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
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"log"
	"time"

	"github.com/alexbeltran/gobacnet"
	"github.com/spf13/cobra"
)

// Flags
var Interface string
var Port int

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
	err = c.WhoIs(1000, 4000)
	log.Println("WHO IS!")
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Duration(5) * time.Second)
	log.Println("DONE")
	c.Close()
}

func init() {
	RootCmd.AddCommand(whoIsCmd)
	whoIsCmd.Flags().StringVarP(&Interface, "interface", "i", "eth0", "Interface e.g. eth0")
	whoIsCmd.Flags().IntVarP(&Port, "port", "p", int(0xBAC0), "Port")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// whoIsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// whoIsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
